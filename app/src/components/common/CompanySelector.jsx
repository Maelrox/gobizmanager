import { useState, useEffect } from 'preact/hooks';
import api from '../../config/axios';
import Select from './Select';

export default function CompanySelector({ value, onChange, className = '' }) {
  const [companies, setCompanies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchCompanies();
  }, []);

  const fetchCompanies = async () => {
    try {
      const response = await api.get('/companies');
      setCompanies(response.data.map(company => ({
        value: company.id,
        label: company.name
      })));
    } catch (err) {
      setError('Failed to fetch companies');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Select
      label="Select Company"
      value={value}
      onChange={onChange}
      options={companies}
      loading={loading}
      error={error}
      className={className}
      placeholder="Select a company"
    />
  );
} 