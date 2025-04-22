import apiClient from '../api/client';

export interface Comment {
  id: string;
  ticketId: string;
  userId: string;
  userName: string;
  userAvatar?: string;
  content: string;
  createdAt: string;
  updatedAt?: string;
}

export interface CommentCreateData {
  ticketId: string;
  content: string;
}

export interface CommentUpdateData {
  content: string;
}

const commentService = {

  async getTicketComments(ticketId: string): Promise<Comment[]> {
    try {
      const response = await apiClient.get(`/tickets/${ticketId}/comments`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching comments for ticket ${ticketId}:`, error);
      throw error;
    }
  },


  async createComment(commentData: CommentCreateData): Promise<Comment> {
    try {
      const response = await apiClient.post(`/tickets/${commentData.ticketId}/comments`, commentData);
      return response.data;
    } catch (error) {
      console.error('Error creating comment:', error);
      throw error;
    }
  },


  async updateComment(commentId: string, ticketId: string, data: CommentUpdateData): Promise<Comment> {
    try {
      const response = await apiClient.put(`/tickets/${ticketId}/comments/${commentId}`, data);
      return response.data;
    } catch (error) {
      console.error(`Error updating comment ${commentId}:`, error);
      throw error;
    }
  },


  async deleteComment(commentId: string, ticketId: string): Promise<void> {
    try {
      await apiClient.delete(`/tickets/${ticketId}/comments/${commentId}`);
    } catch (error) {
      console.error(`Error deleting comment ${commentId}:`, error);
      throw error;
    }
  }
};

export default commentService; 