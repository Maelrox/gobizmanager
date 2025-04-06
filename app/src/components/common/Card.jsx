import { h } from 'preact';

export default function Card({ 
  title, 
  subtitle, 
  children, 
  className = '',
  actions = null 
}) {
  return (
    <div className={`bg-white rounded-lg shadow-sm overflow-hidden ${className}`}>
      {(title || subtitle) && (
        <div className="px-6 py-4 border-b border-neutral-200">
          {title && <h3 className="text-lg font-medium text-neutral-900">{title}</h3>}
          {subtitle && <p className="mt-1 text-sm text-neutral-500">{subtitle}</p>}
        </div>
      )}
      <div className="p-6">
        {children}
      </div>
      {actions && (
        <div className="px-6 py-4 border-t border-neutral-200 bg-neutral-50">
          {actions}
        </div>
      )}
    </div>
  );
} 