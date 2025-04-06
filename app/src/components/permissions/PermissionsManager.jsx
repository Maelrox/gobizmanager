import { useState, useEffect } from 'preact/hooks';
import { h } from 'preact';
import api from '../../config/axios';
import Button from '../common/Button';
import Table from '../common/Table';

const icons = {
  add: (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
    </svg>
  ),
  save: (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
    </svg>
  ),
  delete: (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
    </svg>
  ),
};

export default function PermissionsManager({ companyId }) {
  const [roles, setRoles] = useState([]);
  const [permissions, setPermissions] = useState([]);
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedUser, setSelectedUser] = useState(null);
  const [selectedRole, setSelectedRole] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    permissions: [],
  });

  useEffect(() => {
    fetchData();
  }, [companyId]);

  const fetchData = async () => {
    try {
      const [rolesRes, permissionsRes] = await Promise.all([
        api.get(`/companies/${companyId}/roles`),
        api.get(`/companies/${companyId}/permissions`),
      ]);
      setRoles(rolesRes.data);
      setPermissions(permissionsRes.data);
    } catch (err) {
      setError('Failed to fetch data');
    } finally {
      setLoading(false);
    }
  };

  const searchUsers = async (query) => {
    if (query.length < 2) return;
    try {
      const response = await api.get(`/companies/${companyId}/users/search`, {
        params: { q: query, limit: 50 },
      });
      setUsers(response.data);
    } catch (err) {
      setError('Failed to search users');
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (selectedRole) {
        await api.put(`/companies/${companyId}/roles/${selectedRole.id}`, formData);
      } else {
        await api.post(`/companies/${companyId}/roles`, formData);
      }
      setShowModal(false);
      setFormData({ name: '', description: '', permissions: [] });
      setSelectedRole(null);
      fetchData();
    } catch (err) {
      setError(err.response?.data?.message || 'Operation failed');
    }
  };

  const handleDelete = async (id) => {
    if (confirm('Are you sure you want to delete this role?')) {
      try {
        await api.delete(`/companies/${companyId}/roles/${id}`);
        fetchData();
      } catch (err) {
        setError('Failed to delete role');
      }
    }
  };

  const handleEdit = (role) => {
    setFormData({
      name: role.name,
      description: role.description,
      permissions: role.permissions,
    });
    setSelectedRole(role);
    setShowModal(true);
  };

  const handleAssignRole = async () => {
    if (!selectedUser || !selectedRole) return;
    try {
      await api.post(`/companies/${companyId}/users/${selectedUser.id}/roles`, {
        roleId: selectedRole.id,
      });
      setSelectedUser(null);
      setSelectedRole(null);
      fetchData();
    } catch (err) {
      setError('Failed to assign role');
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  const tableHeaders = ['Name', 'Description', 'Permissions', 'Actions'];

  const renderTableRow = (role) => [
    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-neutral-900">
      {role.name}
    </td>,
    <td className="px-6 py-4 text-sm text-neutral-500">
      {role.description}
    </td>,
    <td className="px-6 py-4 text-sm text-neutral-500">
      {role.permissions.join(', ')}
    </td>,
    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
      <Button
        icon="edit"
        text="Edit"
        onClick={() => handleEdit(role)}
        variant="ghost"
        className="mr-4"
      />
      <Button
        icon="delete"
        text="Delete"
        onClick={() => handleDelete(role.id)}
        variant="ghost"
      />
    </td>,
  ];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-semibold text-neutral-900">Roles & Permissions</h2>
        <Button
          icon="add"
          text="Add Role"
          onClick={() => {
            setFormData({ name: '', description: '', permissions: [] });
            setSelectedRole(null);
            setShowModal(true);
          }}
        />
      </div>

      {error && (
        <div className="bg-accent-50 border-l-4 border-accent-400 p-4 rounded-lg">
          <p className="text-accent-700 text-sm">{error}</p>
        </div>
      )}

      {/* Role Assignment Section */}
      <div className="bg-white rounded-xl shadow-lg p-6">
        <h3 className="text-lg font-medium text-neutral-900 mb-4">Assign Role to User</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-neutral-700 mb-2">
              Search User
            </label>
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => {
                setSearchQuery(e.target.value);
                searchUsers(e.target.value);
              }}
              className="w-full px-4 py-2 rounded-lg border border-neutral-200 focus:border-primary-300 focus:ring-2 focus:ring-primary-200"
              placeholder="Type to search users..."
            />
            {users.length > 0 && (
              <div className="mt-2 max-h-48 overflow-y-auto">
                {users.map((user) => (
                  <div
                    key={user.id}
                    className={`p-2 cursor-pointer hover:bg-neutral-50 ${
                      selectedUser?.id === user.id ? 'bg-primary-50' : ''
                    }`}
                    onClick={() => setSelectedUser(user)}
                  >
                    {user.name} ({user.email})
                  </div>
                ))}
              </div>
            )}
          </div>
          <div>
            <label className="block text-sm font-medium text-neutral-700 mb-2">
              Select Role
            </label>
            <select
              value={selectedRole?.id || ''}
              onChange={(e) => {
                const role = roles.find((r) => r.id === e.target.value);
                setSelectedRole(role);
              }}
              className="w-full px-4 py-2 rounded-lg border border-neutral-200 focus:border-primary-300 focus:ring-2 focus:ring-primary-200"
            >
              <option value="">Select a role</option>
              {roles.map((role) => (
                <option key={role.id} value={role.id}>
                  {role.name}
                </option>
              ))}
            </select>
          </div>
        </div>
        <div className="mt-4 flex justify-end">
          <Button
            text="Assign Role"
            onClick={handleAssignRole}
            disabled={!selectedUser || !selectedRole}
          />
        </div>
      </div>

      {/* Roles Table */}
      <div className="bg-white rounded-xl shadow-lg overflow-hidden">
        <Table
          headers={tableHeaders}
          data={roles}
          renderRow={renderTableRow}
        />
      </div>

      {/* Role Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-xl shadow-lg p-6 max-w-md w-full">
            <h3 className="text-lg font-medium text-neutral-900 mb-4">
              {selectedRole ? 'Edit Role' : 'Add Role'}
            </h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <label htmlFor="name" className="block text-sm font-medium text-neutral-700">
                  Name
                </label>
                <input
                  type="text"
                  id="name"
                  name="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  required
                  className="w-full px-4 py-2 rounded-lg border border-neutral-200 focus:border-primary-300 focus:ring-2 focus:ring-primary-200"
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="description" className="block text-sm font-medium text-neutral-700">
                  Description
                </label>
                <textarea
                  id="description"
                  name="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="w-full px-4 py-2 rounded-lg border border-neutral-200 focus:border-primary-300 focus:ring-2 focus:ring-primary-200"
                  rows="3"
                />
              </div>
              <div className="space-y-2">
                <label className="block text-sm font-medium text-neutral-700">
                  Permissions
                </label>
                <div className="space-y-2">
                  {permissions.map((permission) => (
                    <label key={permission.id} className="flex items-center">
                      <input
                        type="checkbox"
                        checked={formData.permissions.includes(permission.id)}
                        onChange={(e) => {
                          const newPermissions = e.target.checked
                            ? [...formData.permissions, permission.id]
                            : formData.permissions.filter((id) => id !== permission.id);
                          setFormData({ ...formData, permissions: newPermissions });
                        }}
                        className="rounded border-neutral-300 text-primary-600 focus:ring-primary-500"
                      />
                      <span className="ml-2 text-sm text-neutral-700">{permission.name}</span>
                    </label>
                  ))}
                </div>
              </div>
              <div className="flex justify-end space-x-3">
                <Button
                  text="Cancel"
                  onClick={() => setShowModal(false)}
                  variant="ghost"
                />
                <Button
                  text={selectedRole ? 'Update' : 'Create'}
                  type="submit"
                />
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
} 