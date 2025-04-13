import { useState, useEffect } from 'preact/hooks';
import api from '../../config/axios';
import Button from '../common/Button';
import Card from '../common/Card';
import Table from '../common/Table';
import CompanyCard from './CompanyCard';

export default function CompanyList() {
  const [companies, setCompanies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    address: '',
    identifier: '',
    logo: '',
  });
  const [editingId, setEditingId] = useState(null);

  useEffect(() => {
    fetchCompanies();
  }, []);

  const fetchCompanies = async () => {
    try {
      const response = await api.get('/companies');
      setCompanies(response.data);
    } catch (err) {
      setError('Failed to fetch companies');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editingId) {
        await api.put(`/companies/${editingId}`, formData);
      } else {
        await api.post('/companies', formData);
      }
      setShowModal(false);
      setFormData({ name: '', email: '', phone: '', address: '', identifier: '', logo: '' });
      setEditingId(null);
      fetchCompanies();
    } catch (err) {
      setError(err.response?.data?.message || 'Operation failed');
    }
  };

  const handleDelete = async (id) => {
    if (confirm('Are you sure you want to delete this company?')) {
      try {
        await api.delete(`/companies/${id}`);
        fetchCompanies();
      } catch (err) {
        setError('Failed to delete company');
      }
    }
  };

  const handleEdit = (company) => {
    setFormData({
      name: company.name,
      email: company.email,
      phone: company.phone,
      address: company.address,
      identifier: company.identifier,
      logo: company.logo,
    });
    setEditingId(company.id);
    setShowModal(true);
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  const tableHeaders = ['Name', 'Email', 'Phone', 'Identifier', 'Actions'];

  const renderTableRow = (company) => [
    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-neutral-900">
      {company.name}
    </td>,
    <td className="px-6 py-4 whitespace-nowrap text-sm text-neutral-500">
      {company.email}
    </td>,
    <td className="px-6 py-4 whitespace-nowrap text-sm text-neutral-500">
      {company.phone}
    </td>,
    <td className="px-6 py-4 whitespace-nowrap text-sm text-neutral-500">
      {company.identifier}
    </td>,
    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
      <Button
        icon="edit"
        text="Edit"
        onClick={() => handleEdit(company)}
        variant="ghost"
        className="mr-4"
      />
      <Button
        icon="delete"
        text="Delete"
        onClick={() => handleDelete(company.id)}
        variant="ghost"
      />
    </td>,
  ];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-semibold text-neutral-900">Companies</h2>
        <Button
          icon="add"
          text="Add Company"
          onClick={() => {
            setFormData({ name: '', email: '', phone: '', address: '', identifier: '', logo: '' });
            setEditingId(null);
            setShowModal(true);
          }}
        />
      </div>

      {error && (
        <div className="bg-accent-50 border-l-4 border-accent-400 p-4 rounded-lg">
          <p className="text-accent-700 text-sm">{error}</p>
        </div>
      )}

      {/* Desktop Table View */}
      <div className="hidden md:block">
        <div className="bg-white rounded-xl shadow-lg overflow-hidden">
          <Table
            headers={tableHeaders}
            data={companies}
            renderRow={renderTableRow}
          />
        </div>
      </div>

      {/* Mobile Card View */}
      <div className="md:hidden space-y-4">
        {companies.map((company) => (
          <CompanyCard
            key={company.id}
            company={company}
            onEdit={handleEdit}
            onDelete={handleDelete}
          />
        ))}
      </div>

      {/* Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-xl shadow-lg p-6 max-w-md w-full">
            <h3 className="text-lg font-medium text-neutral-900 mb-4">
              {editingId ? 'Edit Company' : 'Add Company'}
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
                  className="w-full px-4 py-3 rounded-lg border border-neutral-200 focus:border-primary-300 focus:ring-2 focus:ring-primary-200 transition-colors duration-200"
                  placeholder="Enter company name"
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="email" className="block text-sm font-medium text-neutral-700">
                  Email
                </label>
                <input
                  type="email"
                  id="email"
                  name="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  required
                  className="w-full px-4 py-3 rounded-lg border border-neutral-200 focus:border-primary-300 focus:ring-2 focus:ring-primary-200 transition-colors duration-200"
                  placeholder="Enter company email"
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="phone" className="block text-sm font-medium text-neutral-700">
                  Phone
                </label>
                <input
                  type="tel"
                  id="phone"
                  name="phone"
                  value={formData.phone}
                  onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                  required
                  className="w-full px-4 py-3 rounded-lg border border-neutral-200 focus:border-primary-300 focus:ring-2 focus:ring-primary-200 transition-colors duration-200"
                  placeholder="Enter company phone"
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="address" className="block text-sm font-medium text-neutral-700">
                  Address
                </label>
                <input
                  type="text"
                  id="address"
                  name="address"
                  value={formData.address}
                  onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                  required
                  className="w-full px-4 py-3 rounded-lg border border-neutral-200 focus:border-primary-300 focus:ring-2 focus:ring-primary-200 transition-colors duration-200"
                  placeholder="Enter company address"
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="identifier" className="block text-sm font-medium text-neutral-700">
                  Identifier
                </label>
                <input
                  type="text"
                  id="identifier"
                  name="identifier"
                  value={formData.identifier}
                  onChange={(e) => setFormData({ ...formData, identifier: e.target.value })}
                  required
                  className="w-full px-4 py-3 rounded-lg border border-neutral-200 focus:border-primary-300 focus:ring-2 focus:ring-primary-200 transition-colors duration-200"
                  placeholder="Enter company identifier"
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="logo" className="block text-sm font-medium text-neutral-700">
                  Logo URL
                </label>
                <input
                  type="text"
                  id="logo"
                  name="logo"
                  value={formData.logo}
                  onChange={(e) => setFormData({ ...formData, logo: e.target.value })}
                  className="w-full px-4 py-3 rounded-lg border border-neutral-200 focus:border-primary-300 focus:ring-2 focus:ring-primary-200 transition-colors duration-200"
                  placeholder="Enter company logo URL (optional)"
                />
              </div>
              <div className="flex justify-end space-x-3">
                <Button
                  icon="cancel"
                  text="Cancel"
                  onClick={() => setShowModal(false)}
                  variant="ghost"
                />
                <Button
                  icon="save"
                  text={editingId ? 'Update' : 'Create'}
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