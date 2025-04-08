import { useState, useEffect } from 'preact/hooks';
import { FaPlus, FaEdit, FaTrash, FaTimes } from 'react-icons/fa';
import Button from '../common/Button';
import CompanySelector from '../common/CompanySelector';
import Table from '../common/Table';
import { userService } from '../../services/userService';

export default function UserManager() {
  const [companyId, setCompanyId] = useState(null);
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    phone: ''
  });

  useEffect(() => {
    if (companyId) {
      console.log('Company ID changed:', companyId);
      fetchUsers();
    }
  }, [companyId]);

  const fetchUsers = async () => {
    setLoading(true);
    setError(null);

    try {
      const usersData = await userService.listUsers(companyId);
      console.log('Received users data:', usersData);
      setUsers(usersData || []);
    } catch (err) {
      console.error('Error fetching users:', err);
      setError(err.response?.data?.message || 'An error occurred');
      setUsers([]);
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteUser = async (id) => {
    if (!confirm('Are you sure you want to delete this user?')) return;

    try {
      await userService.deleteUser(companyId, id);
      setUsers(users.filter(user => user.id !== id));
    } catch (err) {
      setError(err.response?.data?.message || 'An error occurred');
    }
  };

  const handleCompanyChange = (event) => {
    const value = event?.target?.value || event;
    console.log('Company changed to:', value);
    setCompanyId(value);
  };

  const handleEditUser = (user) => {
    console.log('Edit user:', user);
    // TODO: Implement edit functionality
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const newUser = await userService.registerUser(companyId, formData);
      setUsers([...users, newUser]);
      setShowModal(false);
      setFormData({ username: '', password: '', phone: '' });
    } catch (err) {
      setError(err.response?.data?.message || 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const tableHeaders = ['ID', 'Email', 'Actions'];

  const renderTableRow = (user) => [
    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-neutral-900">
      {user.id}
    </td>,
    <td className="px-6 py-4 whitespace-nowrap text-sm text-neutral-500">
      {user.email}
    </td>,
    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
      <div className="flex space-x-2">
        <Button
          icon={<FaEdit className="w-3 h-3" />}
          onClick={() => handleEditUser(user)}
          variant="ghost"
          size="xs"
          className="text-primary-600 hover:text-primary-700"
          title="Edit User"
        />
        <Button
          icon={<FaTrash className="w-3 h-3" />}
          onClick={() => handleDeleteUser(user.id)}
          variant="ghost"
          size="xs"
          className="text-red-600 hover:text-red-700"
          title="Delete User"
        />
      </div>
    </td>
  ];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-semibold text-neutral-900">Users</h1>
        <Button
          icon={<FaPlus className="w-4 h-4" />}
          onClick={() => setShowModal(true)}
          text="Register User"
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
      ) : users.length === 0 ? (
        <div className="text-center py-8 text-neutral-500">
          No users found. Click "Register User" to add a new user.
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          <Table
            headers={tableHeaders}
            data={users}
            renderRow={renderTableRow}
            className="w-full"
          />
        </div>
      )}

      {/* Register User Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl p-6 w-full max-w-md">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-xl font-semibold text-neutral-900">Register User</h2>
              <Button
                icon={<FaTimes className="w-4 h-4" />}
                onClick={() => setShowModal(false)}
                variant="ghost"
                size="sm"
                className="text-neutral-400 hover:text-neutral-500"
              />
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-neutral-700 mb-1">
                  Email
                </label>
                <input
                  type="email"
                  name="username"
                  value={formData.username}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border border-neutral-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-neutral-700 mb-1">
                  Password
                </label>
                <input
                  type="password"
                  name="password"
                  value={formData.password}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border border-neutral-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                  required
                  autoComplete="new-password"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-neutral-700 mb-1">
                  Phone
                </label>
                <input
                  type="tel"
                  name="phone"
                  value={formData.phone}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border border-neutral-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
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
                  onClick={() => setShowModal(false)}
                  variant="ghost"
                  size="sm"
                  className="text-neutral-600 hover:text-neutral-700"
                  text="Cancel"
                />
                <Button
                  type="submit"
                  loading={loading}
                  variant="primary"
                  size="sm"
                  text="Register User"
                />
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
} 