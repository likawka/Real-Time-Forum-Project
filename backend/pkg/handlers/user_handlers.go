package handlers

import (
	"net/http"
	"project-root/pkg/api"
	"project-root/pkg/repositories"
	"project-root/pkg/services"
)

// HandleGetUser retrieves posts or comments by user nickname with pagination and sorting.
// @Summary Get posts or comments by user nickname
// @Description Retrieves posts or comments by user nickname with optional sorting and pagination.
// @Tags users
// @Produce json
// @Param nickname path string true "User Nickname (e.g., admin)"
// @Param type path string true "Type of data to retrieve: posts or comments"
// @Param sort query string false "Sorting criteria: new, old, popular (default: new)"
// @Param pageSize query integer false "Number of items per page (default: 20)"
// @Param page query integer false "Page number (default: 1)"
// @Success 200 {object} api.Response{payload=api.GetUserResponse, pagination=api.GeneralPagination} "Successful operation"
// @Failure 404 {object} api.Response{error=api.ErrorDetails} "User not found"
// @Router /users/{nickname}/{type} [get]
func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user *api.UserResponse

	user, authenticated = services.AuthenticateUser(r)
	params := services.GetRouteParams(r)
	userInfo, err := repositories.GetUserByNickname(params["nickname"])
	if err != nil {
		services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error fetching user", authenticated, user, nil)
		return
	}
	if userInfo == nil {
		services.HTTPError(w, http.StatusNotFound, "User not found", "User not found", authenticated, user, nil)
		return
	}

	sortType, sortBy, page, pageSize := services.ExtractPaginationParams(r, "new", params["type"])

	var (
		payload    api.GetUserResponse
		totalItems int
		totalPages int
	)

	switch params["type"] {
	case "posts":
		postRepo := repositories.NewPostRepository()
		posts, items, pages, err := postRepo.GetPosts(page, pageSize, sortBy, "nickname", params["nickname"])
		if err != nil {
			services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error fetching posts", authenticated, user, nil)
			return
		}
		totalItems = items
		totalPages = pages
		payload = api.GetUserResponse{
			User:  *userInfo,
			Posts: posts,
		}
	case "comments":
		commentRepo := repositories.NewCommentRepository()
		comments, items, pages, err := commentRepo.GetComments(page, pageSize, sortBy, "nickname", params["nickname"], user.ID)
		if err != nil {
			services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error fetching comments", authenticated, user, nil)
			return
		}
		totalItems = items
		totalPages = pages
		payload = api.GetUserResponse{
			User:     *userInfo,
			Comments: *comments,
		}
	default:
		services.HTTPError(w, http.StatusBadRequest, "Bad Request", "Invalid type parameter", authenticated, user, nil)
		return
	}

	if page > totalPages {
		services.HTTPError(w, http.StatusNotFound, "Not Found", "Page not found", authenticated, user, nil)
		return
	}

	pagination := &api.GeneralPagination{
		CurrentPage: page,
		PerPage:     pageSize,
		TotalCount:  totalItems,
		TotalPages:  totalPages,
		OrderBy:     sortType,
	}

	services.RespondWithSuccess(w, http.StatusOK, "Data fetched successfully", authenticated, payload, pagination, user)
}

// HandleGetUsers retrieves all users with pagination.
// @Summary Get all users
// @Description Retrieves all users with optional sorting and pagination.
// @Tags users
// @Produce json
// @Param sort query string false "Sorting criteria: new, old, name_ABC (default: name_ABC)"
// @Param pageSize query integer false "Number of items per page (default: 20)"
// @Param page query integer false "Page number (default: 1)"
// @Success 200 {object} api.Response{payload=api.GetUsersResponse, pagination=api.GeneralPagination} "Successful operation"
// @Router /users [get]
func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user *api.UserResponse

	user, authenticated = services.AuthenticateUser(r)
	sortType, sortBy, page, pageSize := services.ExtractPaginationParams(r, "name_ABC", "users")

	userRepo := repositories.NewUserRepository()
	users, totalItems, totalPages, err := userRepo.GetAllUsers(page, pageSize, sortBy, user.ID)
	if err != nil {
		services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error fetching users", authenticated, user, nil)
		return
	}

	payload := api.GetUsersResponse{
		Users: users,
	}

	pagination := &api.GeneralPagination{
		CurrentPage: page,
		PerPage:     pageSize,
		TotalCount:  totalItems,
		TotalPages:  totalPages,
		OrderBy:     sortType,
	}
	services.RespondWithSuccess(w, http.StatusOK, "Data fetched successfully", authenticated, payload, pagination, user)

}
