import { h } from 'preact';
import { FaPlus, FaEdit, FaTrash, FaTimes, FaCheck } from 'react-icons/fa';

export default function Button({ 
  icon, 
  text, 
  onClick, 
  type = 'button',
  variant = 'neutral',
  className = '',
  disabled = false 
}) {
  const baseStyles = 'inline-flex items-center justify-center px-4 py-2 rounded-lg font-medium transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed';
  
  const variantStyles = {
    neutral: 'bg-neutral-100 text-neutral-700 hover:bg-neutral-200 focus:ring-neutral-500',
    primary: 'bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500',
    danger: 'bg-accent-600 text-white hover:bg-accent-700 focus:ring-accent-500',
    ghost: 'text-neutral-700 hover:bg-neutral-100 focus:ring-neutral-500',
  };

  const iconMap = {
    add: <FaPlus className="w-4 h-4" />,
    edit: <FaEdit className="w-4 h-4" />,
    delete: <FaTrash className="w-4 h-4" />,
    cancel: <FaTimes className="w-4 h-4" />,
    save: <FaCheck className="w-4 h-4" />,
  };

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={`${baseStyles} ${variantStyles[variant]} ${className}`}
    >
      {icon && <span className="mr-2">{iconMap[icon] || icon}</span>}
      {text}
    </button>
  );
} 