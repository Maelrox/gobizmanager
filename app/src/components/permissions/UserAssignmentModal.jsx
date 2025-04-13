import { useState, useEffect } from 'preact/hooks';
import { FaTimes } from 'react-icons/fa';
import Button from '../common/Button';
import Select from '../common/Select';
import { roleService } from '../../services/roleService';

const UserAssignmentModal = ({ role, onClose }) => {
  const [users, setUsers] = useState([]);
  const [selectedUser, setSelectedUser] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const response = await roleService.searchUsers(role.company_id);
      setUsers(response);
    } catch (err) {
      setError('Failed to fetch users');
    }
  };

  const handleAssignUser = async () => {
    if (!selectedUser) {
      setError('Please select a user');
      return;
    }

    setLoading(true);
    try {
      await roleService.assignUserToRole(role.id, selectedUser);
      onClose();
    } catch (err) {
      setError('Failed to assign user');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg w-full max-w-md">
        <div className="flex justify-between items-center p-6 border-b">
          <h2 className="text-2xl font-semibold text-neutral-900">
            Assign User to Role
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

        <div className="p-6">
          <Select
            value={selectedUser}
            onChange={(e) => setSelectedUser(e.target.value)}
            options={users.map(user => ({
              value: user.id,
              label: user.email,
              className: 'text-gray-900'
            }))}
            placeholder="Select a user"
            className="w-full"
          />
        </div>

        <div className="p-6 border-t">
          <div className="flex justify-end space-x-2">
            <Button
              type="button"
              onClick={onClose}
              variant="secondary"
              text="Cancel"
            />
            <Button
              onClick={handleAssignUser}
              disabled={loading}
              text={loading ? 'Assigning...' : 'Assign User'}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default UserAssignmentModal; 