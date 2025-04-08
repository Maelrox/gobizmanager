import { h } from 'preact';

export default function Table({ headers, data, renderRow, columns, className = '' }) {
  // Support both old and new formats
  const tableHeaders = headers || columns?.map(col => col.header);
  const renderTableRow = renderRow || ((row) => 
    columns.map((column, index) => (
      <td
        key={index}
        className="px-6 py-4 whitespace-nowrap text-sm text-neutral-900"
      >
        {column.cell(row)}
      </td>
    ))
  );

  return (
    <div className="overflow-x-auto">
      <table className={`min-w-full divide-y divide-neutral-200 ${className}`}>
        <thead className="bg-neutral-50">
          <tr>
            {tableHeaders.map((header, index) => (
              <th
                key={index}
                className="px-6 py-3 text-left text-xs font-medium text-neutral-500 uppercase tracking-wider"
              >
                {header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-neutral-200">
          {data.map((row, rowIndex) => (
            <tr key={rowIndex} className="hover:bg-neutral-50">
              {renderTableRow(row)}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
} 