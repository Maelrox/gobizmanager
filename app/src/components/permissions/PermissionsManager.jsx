import { useState, useEffect } from 'preact/hooks';
import { FaEdit, FaUserPlus, FaKey, FaTrash, FaPlus } from 'react-icons/fa';
import Button from '../common/Button';
import Table from '../common/Table';
import CompanySelector from '../common/CompanySelector';
import Select from '../common/Select';
import Card from '../common/Card';
import RoleCard from './RoleCard';
import RoleForm from './RoleForm';
import { roleService } from '../../services/roleService';
import PermissionModal from './PermissionModal';
import UserAssignmentModal from './UserAssignmentModal';

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

export default function PermissionsManager() {
  const [companyId, setCompanyId] = useState(null);
  const [roles, setRoles] = useState([]);
  const [permissions, setPermissions] = useState([]);
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showForm, setShowForm] = useState(false);
  const [selectedRole, setSelectedRole] = useState(null);
  const [showUserAssignment, setShowUserAssignment] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    description: ''
  });

  useEffect(() => {
    if (companyId) {
      fetchData();
    }
  }, [companyId]);

  const fetchData = async () => {
    setLoading(true);
    setError(null);

    try {
      const [rolesData, permissionsData] = await Promise.all([
        roleService.listRoles(companyId),
        roleService.getPermissions(companyId)
      ]);

      const rolesWithPermissions = rolesData.map(role => ({
        ...role,
        permissions: role.permissions || []
      }));

      setRoles(rolesWithPermissions);
      setPermissions(permissionsData);
    } catch (err) {
      setError(err.response?.data?.message || 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const searchUsers = async (query) => {
    if (query.length < 2) return;
    try {
      const response = await api.get(`/companies/${companyId}/users/search`, {
        params: { q: query, limit: 50 }
      });
      setUsers(response.data);
    } catch (err) {
      setError('Failed to search users');
    }
  };

  const handleCreateRole = () => {
    setSelectedRole(null);
    setShowForm(true);
  };

  const handleEditRole = (role) => {
    setSelectedRole(role);
  };

  const handleCloseModal = () => {
    setSelectedRole(null);
  };

  const handleUpdatePermissions = () => {
    fetchData();
  };

  const handleDeleteRole = async (id) => {
    if (!confirm('Are you sure you want to delete this role?')) return;

    try {
      await roleService.deleteRole(id);
      setRoles(roles.filter(role => role.id !== id));
    } catch (err) {
      setError(err.response?.data?.message || 'An error occurred');
    }
  };

  const handleSubmitRole = async (formData) => {
    try {
      if (selectedRole) {
        const updatedRole = await roleService.updateRole(selectedRole.id, formData);
        setRoles(roles.map(role => 
          role.id === selectedRole.id ? { ...role, ...updatedRole } : role
        ));
      } else {
        const newRole = await roleService.createRole(formData);
        setRoles([...roles, newRole]);
      }
    } catch (err) {
      throw err;
    }
  };

  const handleAddUser = (role) => {
    setSelectedRole(role);
    setShowUserAssignment(true);
  };

  const handleCloseUserAssignment = () => {
    setShowUserAssignment(false);
    setSelectedRole(null);
  };

  const handleCompanyChange = (event) => {
    // Extract the value from the event object if it exists
    const value = event?.target?.value || event;
    setCompanyId(value);
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-semibold text-neutral-900">Roles & Permissions</h1>
        <Button
          icon={<FaPlus className="w-4 h-4" />}
          onClick={handleCreateRole}
          text='Create role'
        />
      </div>

      <CompanySelector
        value={companyId}
        onChange={handleCompanyChange}
        className="w-full max-w-md"
      />

      {loading ? (
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600" />
        </div>
      ) : error ? (
        <div className="text-red-600">{error}</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {roles.map(role => (
            <RoleCard
              key={role.id}
              role={role}
              onEdit={handleEditRole}
              onDelete={handleDeleteRole}
              onAddPermission={() => handleEditRole(role)}
              onAddUser={() => handleAddUser(role)}
            />
          ))}
        </div>
      )}

      {showForm && (
        <RoleForm
          role={selectedRole}
          permissions={permissions}
          companyId={companyId}
          onClose={() => {
            setShowForm(false);
            setSelectedRole(null);
          }}
          onSubmit={handleSubmitRole}
        />
      )}

      {selectedRole && !showUserAssignment && (
        <PermissionModal
          role={selectedRole}
          companyId={companyId}
          onClose={handleCloseModal}
          onUpdate={handleUpdatePermissions}
        />
      )}

      {selectedRole && showUserAssignment && (
        <UserAssignmentModal
          role={selectedRole}
          onClose={handleCloseUserAssignment}
        />
      )}
    </div>
  );
} 