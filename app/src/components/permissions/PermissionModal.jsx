import React, { useState, useEffect } from 'react';
import { FaTimes, FaKey, FaCog, FaPlus } from 'react-icons/fa';
import Button from '../common/Button';
import { roleService } from '../../services/roleService';
import Input from '../common/Input';

const PermissionModal = ({
  role,
  companyId,
  onClose,
  onUpdate
}) => {
  const [permissions, setPermissions] = useState([]);
  const [selectedPermissions, setSelectedPermissions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState('permissions');
  const [selectedPermission, setSelectedPermission] = useState(null);
  const [moduleActions, setModuleActions] = useState([]);
  const [selectedModuleActions, setSelectedModuleActions] = useState([]);
  const [newPermission, setNewPermission] = useState({
    name: '',
    description: ''
  });

  useEffect(() => {
    if (role && companyId) {
      loadPermissions();
    }
  }, [role, companyId]);

  const loadPermissions = async () => {
    try {
      setLoading(true);
      const data = await roleService.getPermissions(companyId);
      setPermissions(data);
      // Initialize selected permissions from role's current permissions
      if (role && role.permissions) {
        setSelectedPermissions(role.permissions.map(p => p.id));
      }
    } catch (err) {
      setError('Failed to load permissions');
    } finally {
      setLoading(false);
    }
  };

  const handlePermissionChange = (permissionId) => {
    setSelectedPermissions(prev => {
      if (prev.includes(permissionId)) {
        return prev.filter(id => id !== permissionId);
      } else {
        return [...prev, permissionId];
      }
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      setLoading(true);
      await roleService.updateRolePermissions(role.id, selectedPermissions);
      onUpdate();
      onClose();
    } catch (err) {
      setError('Failed to update permissions');
    } finally {
      setLoading(false);
    }
  };

  const handleCreatePermission = async (e) => {
    e.preventDefault();
    try {
      setLoading(true);
      await roleService.createPermission({
        company_id: parseInt(companyId),
        name: newPermission.name,
        description: newPermission.description,
        role_id: parseInt(role.id)
      });
      // Refresh permissions list
      await loadPermissions();
      setNewPermission({ name: '', description: '' });
    } catch (err) {
      setError('Failed to create permission');
    } finally {
      setLoading(false);
    }
  };

  const loadModuleActions = async () => {
    try {
      setLoading(true);
      const [actions, permissionActions] = await Promise.all([
        roleService.getModuleActions(),
        roleService.getPermissionModuleActions(selectedPermission.id)
      ]);
      
      // Group actions by module
      const groupedActions = actions.reduce((acc, action) => {
        if (!acc[action.module_name]) {
          acc[action.module_name] = [];
        }
        acc[action.module_name].push(action);
        return acc;
      }, {});

      setModuleActions(groupedActions);
      setSelectedModuleActions(permissionActions.map(pa => pa.module_action_id));
    } catch (err) {
      setError('Failed to load module actions');
    } finally {
      setLoading(false);
    }
  };

  const handleModuleActionChange = (moduleActionId) => {
    setSelectedModuleActions(prev => {
      if (prev.includes(moduleActionId)) {
        return prev.filter(id => id !== moduleActionId);
      } else {
        return [...prev, moduleActionId];
      }
    });
  };

  const handleUpdateModuleActions = async (e) => {
    e.preventDefault();
    try {
      setLoading(true);
      await roleService.updatePermissionModuleActions(selectedPermission.id, selectedModuleActions);
      setShowModuleActionsModal(false);
      await loadPermissions();
    } catch (err) {
      setError('Failed to update module actions');
    } finally {
      setLoading(false);
    }
  };

  const tabs = [
    { 
      id: 'permissions', 
      label: 'Permissions', 
      icon: <FaKey className="w-4 h-4" />,
      description: 'Manage role permissions by selecting or deselecting available permissions. These permissions define what actions users with this role can perform.'
    },
    { 
      id: 'module-actions', 
      label: 'Module Actions', 
      icon: <FaCog className="w-4 h-4" />,
      description: 'Configure specific module actions for each permission. This allows fine-grained control over what operations are allowed within each module.'
    },
    { 
      id: 'create', 
      label: 'Create Permission', 
      icon: <FaPlus className="w-4 h-4" />,
      description: 'Create a new permission with a unique name and description. This permission can then be assigned to roles and configured with specific module actions.'
    }
  ];

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg w-full max-w-4xl h-[80vh] flex flex-col">
        <div className="flex justify-between items-center p-6 border-b">
          <h2 className="text-2xl font-semibold text-neutral-900">
            Manage Permissions
          </h2>
          <Button
            icon={<FaTimes className="w-4 h-4" />}
            onClick={onClose}
            variant="ghost"
            size="sm"
            className="text-neutral-400 hover:text-neutral-500"
          />
        </div>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 mx-6 my-4 rounded">
            {error}
          </div>
        )}

        <div className="border-b border-gray-200">
          <nav className="flex space-x-8 px-6" aria-label="Tabs">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`
                  flex items-center space-x-2 py-4 px-1 border-b-2 font-medium text-sm transition-colors duration-200
                  ${activeTab === tab.id
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }
                `}
              >
                {tab.icon}
                <span>{tab.label}</span>
              </button>
            ))}
          </nav>
        </div>

        <div className="flex-1 flex flex-col min-h-0">
          <div className="p-6">
            <p className="text-gray-600 text-sm mb-4">
              {tabs.find(tab => tab.id === activeTab)?.description}
            </p>
          </div>

          <div className="flex-1 overflow-y-auto px-6">
            {activeTab === 'permissions' && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {permissions.map(permission => (
                  <div key={permission.id} className="flex items-center p-3 bg-gray-50 rounded-lg">
                    <input
                      type="checkbox"
                      id={`permission-${permission.id}`}
                      checked={selectedPermissions.includes(permission.id)}
                      onChange={() => handlePermissionChange(permission.id)}
                      className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                    />
                    <label htmlFor={`permission-${permission.id}`} className="ml-3 block text-sm font-medium text-gray-700">
                      {permission.name}
                    </label>
                  </div>
                ))}
              </div>
            )}

            {activeTab === 'module-actions' && (
              <div className="space-y-6">
                <div className="flex items-center space-x-4">
                  <select
                    value={selectedPermission?.id || ''}
                    onChange={(e) => {
                      const perm = permissions.find(p => p.id === parseInt(e.target.value));
                      setSelectedPermission(perm);
                      if (perm) loadModuleActions();
                    }}
                    className="block w-64 rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                  >
                    <option value="">Select a permission</option>
                    {permissions.map(p => (
                      <option key={p.id} value={p.id}>{p.name}</option>
                    ))}
                  </select>
                </div>

                {selectedPermission && (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {Object.entries(moduleActions).map(([moduleName, actions]) => (
                      <div key={moduleName} className="bg-gray-50 rounded-lg p-4">
                        <h3 className="font-medium text-gray-900 mb-3">{moduleName}</h3>
                        <div className="space-y-2">
                          {actions.map(action => (
                            <div key={action.id} className="flex items-center">
                              <input
                                type="checkbox"
                                id={`action-${action.id}`}
                                checked={selectedModuleActions.includes(action.id)}
                                onChange={() => handleModuleActionChange(action.id)}
                                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                              />
                              <label htmlFor={`action-${action.id}`} className="ml-3 block text-sm text-gray-700">
                                {action.name}
                              </label>
                            </div>
                          ))}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}

            {activeTab === 'create' && (
              <div className="space-y-6 max-w-lg mx-auto">
                <div>
                  <label htmlFor="permission-name" className="block text-sm font-medium text-gray-700 mb-1">
                    Permission Name
                  </label>
                  <Input
                    id="permission-name"
                    name="name"
                    value={newPermission.name}
                    onChange={(e) => setNewPermission(prev => ({ ...prev, name: e.target.value }))}
                    placeholder="e.g., manage_users"
                    required
                  />
                </div>
                <div>
                  <label htmlFor="permission-description" className="block text-sm font-medium text-gray-700 mb-1">
                    Description
                  </label>
                  <Input
                    id="permission-description"
                    name="description"
                    value={newPermission.description}
                    onChange={(e) => setNewPermission(prev => ({ ...prev, description: e.target.value }))}
                    placeholder="e.g., Allows managing user accounts"
                    required
                  />
                </div>
              </div>
            )}
          </div>

          <div className="p-6 border-t mt-auto">
            <div className="flex justify-end space-x-2">
              {activeTab === 'permissions' && (
                <>
                  <Button
                    type="button"
                    onClick={onClose}
                    variant="secondary"
                    text="Cancel"
                  />
                  <Button
                    onClick={handleSubmit}
                    disabled={loading}
                    text={loading ? 'Updating...' : 'Update Permissions'}
                  />
                </>
              )}

              {activeTab === 'module-actions' && (
                <>
                  <Button
                    type="button"
                    onClick={onClose}
                    variant="secondary"
                    text="Cancel"
                  />
                  <Button
                    onClick={handleUpdateModuleActions}
                    disabled={loading || !selectedPermission}
                    text={loading ? 'Updating...' : 'Update Actions'}
                  />
                </>
              )}

              {activeTab === 'create' && (
                <>
                  <Button
                    type="button"
                    onClick={onClose}
                    variant="secondary"
                    text="Cancel"
                  />
                  <Button
                    onClick={handleCreatePermission}
                    disabled={loading}
                    text={loading ? 'Creating...' : 'Create Permission'}
                  />
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PermissionModal; 