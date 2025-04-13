import api from '../config/axios';

export const roleService = {
  async listRoles(companyId) {
    // If companyId is an event object, extract the value
    const id = companyId?.target?.value || companyId;
    if (!id) {
      throw new Error('Company ID is required');
    }
    const response = await api.get(`/rbac/roles/company/${id}`);
    return response.data;
  },

  async getPermissions(companyId) {
    // If companyId is an event object, extract the value
    const id = companyId?.target?.value || companyId;
    if (!id) {
      throw new Error('Company ID is required');
    }
    const response = await api.get(`/rbac/permissions/company/${id}`);
    return response.data;
  },

  async createRole(roleData) {
    const response = await api.post('/rbac/roles', roleData);
    return response.data;
  },

  async updateRole(roleId, roleData) {
    const response = await api.put(`/rbac/roles/${roleId}`, roleData);
    return response.data;
  },

  async deleteRole(roleId) {
    await api.delete(`/rbac/roles/${roleId}`);
  },

  async searchUsers(companyId) {
    const response = await api.get(`/users/search?companyId=${companyId}`);
    return response.data;
  },

  async updateRolePermissions(roleId, permissionIds) {
    const response = await api.put(`/rbac/roles/${roleId}/permissions`, {
      permission_ids: permissionIds
    });
    return response.data;
  },

  async assignUserToRole(roleId, userId) {
    const response = await api.post('/rbac/roles/assign', {
      roleId,
      userId
    });
    return response.data;
  },

  createPermission: async (permissionData) => {
    const response = await api.post('/rbac/permissions', permissionData);
    return response.data;
  },

  getModuleActions: async () => {
    const response = await api.get('/rbac/module-actions');
    return response.data;
  },

  getPermissionModuleActions: async (permissionId) => {
    const response = await api.get(`/rbac/permissions/${permissionId}/module-actions`);
    return response.data;
  },

  updatePermissionModuleActions: async (permissionId, moduleActionIds) => {
    const response = await api.put(`/rbac/permissions/${permissionId}/module-actions`, {
      module_action_ids: moduleActionIds
    });
    return response.data;
  }
}; 