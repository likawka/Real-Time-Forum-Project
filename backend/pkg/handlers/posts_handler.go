package handlers

import (
	"log"
	"net/http"

	"encoding/json"
	"project-root/pkg/api"
	"project-root/pkg/repositories"
	"project-root/pkg/services"
	"strconv"
)

// HandleGetPosts retrieves all posts from the database, with optional sorting and pagination.
// @Summary Get all posts
// @Description Retrieves all posts from the database, with optional sorting and pagination.
// @Tags posts
// @Produce json
// @Param sort query string false "Sorting criteria: new, old, popular (default: new)"
// @Param pageSize query integer false "Number of items per page (default: 20)"
// @Param page query integer false "Page number (default: 1)"
// @Success 200 {object} api.Response{payload=api.PostsResponse, pagination=api.GeneralPagination} "Successful operation"
// @Failure 404 {object} api.Response{error=api.ErrorDetails} "Posts not found"
// @Router /posts [get]
func HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user *api.UserResponse

	user, authenticated = services.AuthenticateUser(r)

	sortType, sortBy, page, pageSize := services.ExtractPaginationParams(r, "new", "posts")

	postRepo := repositories.NewPostRepository()

	posts, totalItems, totalPages, err := postRepo.GetPosts(page, pageSize, sortBy, "", "")
	if err != nil {
		log.Println("Error fetching posts:", err)
		services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error fetching posts", authenticated, user, nil)
		return
	}

	if page > totalPages {
		services.HTTPError(w, http.StatusNotFound, "Not Found", "Page not found", authenticated, user, nil)
		return
	}

	payload := api.PostsResponse{
		Posts: posts,
	}
	pagination := &api.GeneralPagination{
		CurrentPage: page,
		PerPage:     pageSize,
		TotalCount:  totalItems,
		TotalPages:  totalPages,
		OrderBy:     sortType,
	}

	services.RespondWithSuccess(w, http.StatusOK, "Posts fetched successfully", authenticated, payload, pagination, user)
}

// HandleGetPostAndComments retrieves a specific post and its comments with pagination and sorting.
// @Summary Get a specific post and its comments
// @Description Retrieves a specific post and its comments by ID, with optional sorting and pagination for comments.
// @Tags posts
// @Produce json
// @Param id path integer true "Post ID"
// @Param sort query string false "Sorting criteria: new, old, popular (default: new)"
// @Param pageSize query integer false "Number of items per page (default: 20)"
// @Param page query integer false "Page number (default: 1)"
// @Success 200 {object} api.Response{payload=api.PostAndCommentsResponse, pagination=api.GeneralPagination} "Successful operation"
// @Failure 400 {object} api.Response{error=api.ErrorDetails} "Invalid post ID"
// @Router /posts/{id} [get]
func HandleGetPostAndComments(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user *api.UserResponse

	user, authenticated = services.AuthenticateUser(r)

	params := services.GetRouteParams(r)
	postID, err := strconv.Atoi(params["postId"])
	if err != nil {
		services.HTTPError(w, http.StatusBadRequest, "Bad Request", "Invalid post ID", authenticated, user, nil)
		return
	}

	sortType, sortBy, page, pageSize := services.ExtractPaginationParams(r, "popular", "comments")
	postRepo := repositories.NewPostRepository()
	post, err := postRepo.GetPostByID(postID, user.ID)
	if err != nil {
		services.HTTPError(w, http.StatusNotFound, "Not Found", "Post not found", authenticated, user, nil)
		return
	}

	commentRepo := repositories.NewCommentRepository()
	comments, totalItems, totalPages, err := commentRepo.GetComments(page, pageSize, sortBy, "postID", postID, user.ID)
	if err != nil {
		services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error fetching comments", authenticated, user, nil)
		return
	}

	payload := api.PostAndCommentsResponse{
		Post:     *post,
		Comments: comments,
	}

	pagination := &api.GeneralPagination{
		CurrentPage: page,
		PerPage:     pageSize,
		TotalCount:  totalItems,
		TotalPages:  totalPages,
		OrderBy:     sortType,
	}

	services.RespondWithSuccess(w, http.StatusOK, "Post and comments fetched successfully", authenticated, payload, pagination, user)
}

// HandleCreatePost creates a new post
// @Summary Create a new post
// @Description Creates a new post with the provided data.
// @Tags posts
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param post body api.PostCreateRequest true "Post data to create"
// @Success 201 {object} api.Response{payload=api.PostCreateResponse} "Post created successfully"
// @Failure 400 {object} api.Response{error=api.ErrorDetails} "Invalid request payload"
// @Failure 401 {object} api.Response{error=api.ErrorDetails} "Unauthorized: User is not authenticated"
// @Failure 500 {object} api.Response{error=api.ErrorDetails} "Internal Server Error: Error creating post"
// @Router /posts [post]
func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user *api.UserResponse

	user, authenticated = services.AuthenticateUser(r)
	if !authenticated {
		services.HTTPError(w, http.StatusUnauthorized, "Unauthorized", "User is not authenticated", false, nil, nil)
		return
	}

	var postForm api.PostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&postForm); err != nil {
		services.HTTPError(w, http.StatusBadRequest, "Bad Request", "Invalid request payload", authenticated, user, nil)
		return
	}

	postForm.Title = services.TrimAndNormalizeSpaces(postForm.Title)
	postForm.Content = services.TrimAndNormalizeSpaces(postForm.Content)

	validationErrors := services.ValidateOperation("post", postForm)
	if len(validationErrors) > 0 {
		services.HTTPError(w, http.StatusBadRequest, "Validation error", "Validation error", authenticated, nil, validationErrors)
		return
	}
	postForm.UserID = user.ID

	postRepo := repositories.NewPostRepository()

	if err := postRepo.Create(&postForm); err != nil {
		log.Println("Error creating post:", err)
		services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error creating post", authenticated, user, nil)
		return
	}

	if err := postRepo.AddPostCategories(postForm.ID, services.ParseCategories(postForm.Categories)); err != nil {
		log.Println("Error creating categories:", err)
		services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error creating categories", authenticated, user, nil)
		return
	}

	services.RespondWithJSON(w, http.StatusCreated, api.Response{
		Status:        "success",
		Message:       "Post created successfully",
		Payload:       api.PostCreateResponse{ID: postForm.ID},
		Authenticated: authenticated,
		User:          user,
	})
}

// HandleCreateComment creates a new comment for a post
// @Summary Create a new comment for a post
// @Description Creates a new comment for a specified post ID. Requires authentication.
// @Tags posts
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param postId path integer true "Post ID"
// @Param body body api.CommentCreateRequest true "Comment data"
// @Success 201 {object} api.Response{payload=api.CommentResponse} "Comment created successfully"
// @Failure 400 {object} api.Response{error=api.ErrorDetails} "Bad request"
// @Failure 401 {object} api.Response{error=api.ErrorDetails} "Unauthorized"
// @Failure 500 {object} api.Response{error=api.ErrorDetails} "Internal server error"
// @Router /posts/{postId}/comments [post]
func HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user *api.UserResponse

	user, authenticated = services.AuthenticateUser(r)
	if !authenticated {
		services.HTTPError(w, http.StatusUnauthorized, "Unauthorized", "User is not authenticated", false, nil, nil)
		return
	}

	params := services.GetRouteParams(r)
	postID, err := strconv.Atoi(params["postId"])
	if err != nil {
		services.HTTPError(w, http.StatusBadRequest, "Bad Request", "Invalid post ID", authenticated, user, nil)
		return
	}

	var commentForm api.CommentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&commentForm); err != nil {
		services.HTTPError(w, http.StatusBadRequest, "Bad Request", "Invalid request payload", authenticated, user, nil)
		return
	}

	commentForm.Content = services.TrimAndNormalizeSpaces(commentForm.Content)

	validationErrors := services.ValidateOperation("comment", commentForm)
	if len(validationErrors) > 0 {
		services.HTTPError(w, http.StatusBadRequest, "Validation error", "Validation error", authenticated, nil, validationErrors)
		return
	}
	commentForm.PostID = postID
	commentForm.UserID = user.ID
	commentRepo := repositories.NewCommentRepository()

	if err := commentRepo.Create(&commentForm); err != nil {
		log.Println("Error creating comment:", err)
		services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error creating comment", authenticated, user, nil)
		return
	}
	payload := api.CommentResponse{
		Comment: api.Comment{
			ID:        commentForm.ID,
			PostID:    commentForm.PostID,
			UserID:    commentForm.UserID,
			Nickname:  user.Nickname,
			Content:   commentForm.Content,
			CreatedAt: commentForm.CreatedAt,
			Rate:      api.Rate{Rate: 0, Status: ""},
		}}

	services.RespondWithJSON(w, http.StatusCreated, api.Response{
		Status:        "success",
		Message:       "Comment added successfully",
		Payload:       payload,
		Authenticated: authenticated,
		User:          user,
	})
}
