import { h } from 'preact';
import Button from '../common/Button';

export default function CompanyCard({ company, onEdit, onDelete }) {
  return (
    <div className="bg-white rounded-xl shadow-lg overflow-hidden">
      <div className="p-4">
        <h3 className="text-lg font-medium text-neutral-900 mb-2">{company.name}</h3>
        <div className="space-y-1 text-sm text-neutral-500">
          <p>{company.email}</p>
          <p>{company.phone}</p>
          <p>{company.identifier}</p>
        </div>
      </div>
      <div className="bg-neutral-50 px-4 py-3 flex justify-end space-x-3">
        <Button
          icon="edit"
          text="Edit"
          onClick={() => onEdit(company)}
          variant="ghost"
        />
        <Button
          icon="delete"
          text="Delete"
          onClick={() => onDelete(company.id)}
          variant="ghost"
        />
      </div>
    </div>
  );
} 