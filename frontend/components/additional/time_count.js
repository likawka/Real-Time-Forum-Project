export function timeAgo(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diff = now - date;
  
    const years = Math.floor(diff / (365 * 24 * 60 * 60 * 1000));
    const months = Math.floor((diff % (365 * 24 * 60 * 60 * 1000)) / (30 * 24 * 60 * 60 * 1000));
    const days = Math.floor((diff % (30 * 24 * 60 * 60 * 1000)) / (24 * 60 * 60 * 1000));
    const hours = Math.floor((diff % (24 * 60 * 60 * 1000)) / (60 * 60 * 1000));
    const minutes = Math.floor((diff % (60 * 60 * 1000)) / (60 * 1000));
    const seconds = Math.floor((diff % (60 * 1000)) / 1000);
  
    let result = '';
  
    if (years > 0) {
      result += `${years} year${years > 1 ? 's' : ''}`;
      if (months > 0) {
        result += `, ${months} month${months > 1 ? 's' : ''}`;
      }
    } else if (months > 0) {
      result += `${months} month${months > 1 ? 's' : ''}`;
      if (days > 0) {
        result += `, ${days} day${days > 1 ? 's' : ''}`;
      }
    } else if (days > 0) {
      result += `${days} day${days > 1 ? 's' : ''}`;
      if (hours > 0) {
        result += `, ${hours} hour${hours > 1 ? 's' : ''}`;
      }
    } else if (hours > 0) {
      result += `${hours} hour${hours > 1 ? 's' : ''}`;
      if (minutes > 0) {
        result += `, ${minutes} minute${minutes > 1 ? 's' : ''}`;
      }
    } else if (minutes > 0) {
      result += `${minutes} minute${minutes > 1 ? 's' : ''}`;
      if (seconds > 0) {
        result += `, ${seconds} second${seconds > 1 ? 's' : ''}`;
      }
    } else if (seconds > 0) {
      result += `${seconds} second${seconds > 1 ? 's' : ''}`;
    }
  
    if (result !== '') {
      result += ' ago';
    } else {
      result = 'just now';
    }
  
    return result;
  }

