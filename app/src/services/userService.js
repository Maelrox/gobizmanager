import api from '../config/axios';

export const userService = {
  async listUsers(companyId) {
    // If companyId is an event object, extract the value
    const id = companyId?.target?.value || companyId;
    if (!id) {
      throw new Error('Company ID is required');
    }
    const response = await api.get(`/users/search?companyId=${id}`);
    return response.data || [];
  },

  async searchUsers(companyId, query) {
    // If companyId is an event object, extract the value
    const id = companyId?.target?.value || companyId;
    if (!id) {
      throw new Error('Company ID is required');
    }
    const response = await api.get(`/users/search`, {
      params: { companyId: id, q: query }
    });
    return response.data || [];
  },

  async registerUser(companyId, userData) {
    // If companyId is an event object, extract the value
    const id = companyId?.target?.value || companyId;
    if (!id) {
      throw new Error('Company ID is required');
    }
    const response = await api.post(`/company-users/register`, { ...userData, company_id: parseInt(id) });
    return response.data;
  },

  async updateUser(companyId, userId, userData) {
    // If companyId is an event object, extract the value
    const id = companyId?.target?.value || companyId;
    if (!id) {
      throw new Error('Company ID is required');
    }
    const response = await api.put(`/users/${userId}`, { ...userData, companyId: id });
    return response.data;
  },

  async deleteUser(companyId, userId) {
    // If companyId is an event object, extract the value
    const id = companyId?.target?.value || companyId;
    if (!id) {
      throw new Error('Company ID is required');
    }
    await api.delete(`/users/${userId}?companyId=${id}`);
  }
}; 