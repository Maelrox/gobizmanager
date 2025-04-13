import { h } from 'preact';

export default function Select({
  label,
  value,
  onChange,
  options = [],
  loading = false,
  error,
  className = '',
  disabled = false,
  placeholder = 'Select an option'
}) {
  return (
    <div className={`space-y-2 ${className}`}>
      {label && (
        <label className="block text-sm font-medium text-neutral-700">
          {label}
        </label>
      )}
      <select
        value={value}
        onChange={onChange}
        disabled={disabled || loading}
        className={`w-full px-4 py-2 rounded-lg border ${
          error ? 'border-accent-500' : 'border-neutral-200'
        } focus:border-primary-300 focus:ring-2 focus:ring-primary-200 transition-colors duration-200 ${
          disabled || loading ? 'bg-neutral-50 cursor-not-allowed' : 'bg-white'
        }`}
      >
        <option value="">{placeholder}</option>
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
      {loading && (
        <div className="absolute right-3 top-1/2 -translate-y-1/2">
          <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary-600"></div>
        </div>
      )}
      {error && (
        <p className="text-sm text-accent-600">{error}</p>
      )}
    </div>
  );
} 