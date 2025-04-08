import { useState, useEffect } from 'preact/hooks';
import { FaTimes } from 'react-icons/fa';
import Button from '../common/Button';
import Input from '../common/Input';
import { roleService } from '../../services/roleService';

export default function RoleForm({ 
  role, 
  companyId,
  onClose, 
  onSubmit 
}) {
  const [formData, setFormData] = useState({
    name: '',
    description: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (role) {
      setFormData({
        name: role.name || '',
        description: role.description || ''
      });
    }
  }, [role]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const roleData = {
        ...formData,
        company_id: Number(companyId)
      };
      await roleService.createRole(roleData);
      onClose();
    } catch (err) {
      setError(err.response?.data?.message || 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl p-6 w-full max-w-md">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-xl font-semibold text-neutral-900">
            {role ? 'Edit Role' : 'Create Role'}
          </h2>
          <Button
            icon={<FaTimes className="w-4 h-4" />}
            onClick={onClose}
            variant="ghost"
            size="sm"
            className="text-neutral-400 hover:text-neutral-500"
          />
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-neutral-700 mb-1">
              Name
            </label>
            <Input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Enter role name"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-neutral-700 mb-1">
              Description
            </label>
            <Input
              type="text"
              name="description"
              value={formData.description}
              onChange={handleChange}
              placeholder="Enter role description"
              required
            />
          </div>

          {error && (
            <div className="text-sm text-red-600">
              {error}
            </div>
          )}

          <div className="flex justify-end space-x-3">
            <Button
              type="button"
              onClick={onClose}
              variant="ghost"
              text="Cancel"
            />
            <Button
              type="submit"
              loading={loading}
              text={role ? 'Update Role' : 'Create Role'}
            />
          </div>
        </form>
      </div>
    </div>
  );
} 