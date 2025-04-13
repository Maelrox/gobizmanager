import { FaEdit, FaKey, FaUserPlus, FaTrash } from 'react-icons/fa';
import Button from '../common/Button';

export default function RoleCard({
  role,
  onEdit,
  onDelete,
  onAddPermission,
  onAddUser
}) {
  return (
    <div className="p-6 flex flex-col h-full border border-neutral-200 hover:border-primary-300 hover:shadow-lg transition-all duration-200">
      <div className="flex-1">
        <div className="flex justify-between items-start mb-4">
          <h3 className="text-lg font-semibold text-neutral-900">{role.name}</h3>
        </div>
        <p className="text-sm text-neutral-500 mb-6">{role.description}</p>

        <div className="space-y-6">
          <div className="space-y-2">
            <h4 className="text-sm font-medium text-neutral-700">Permissions</h4>
            <div className="flex flex-wrap gap-2">
              {role.permissions && role.permissions.length > 0 ? (
                role.permissions.map((permission) => (
                  <span
                    key={permission.id}
                    className="px-2.5 py-1 text-xs font-medium bg-primary-50 text-primary-700 border border-primary-100"
                  >
                    {permission.name}
                  </span>
                ))
              ) : (
                <span className="text-sm text-neutral-400">No permissions assigned</span>
              )}
            </div>
          </div>

          <div className="space-y-2">
            <h4 className="text-sm font-medium text-neutral-700">Users</h4>
            <div className="flex flex-wrap gap-2">
              {role.users && role.users.length > 0 ? (
                role.users.map((user) => (
                  <span
                    key={user.id}
                    className="px-2.5 py-1 text-xs font-medium bg-neutral-50 text-neutral-700 rounded-full border border-neutral-200"
                  >
                    {user.name}
                  </span>
                ))
              ) : (
                <span className="text-sm text-neutral-400">No users assigned</span>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="mt-6 pt-4 border-t border-neutral-100 flex justify-center items-center">
        <div className="flex space-x-1">
          <Button
            icon={<FaEdit className="w-3 h-3" />}
            onClick={() => onEdit(role)}
            variant="ghost"
            size="xs"
            className="text-primary-600 hover:text-primary-700 p-1"
            title="Edit Role"
          />
          <Button
            icon={<FaKey className="w-3 h-3" />}
            onClick={() => onAddPermission(role)}
            variant="ghost"
            size="xs"
            className="text-primary-600 hover:text-primary-700 p-1"
            title="Add Permission"
          />
          <Button
            icon={<FaUserPlus className="w-3 h-3" />}
            onClick={() => onAddUser(role)}
            variant="ghost"
            size="xs"
            className="text-primary-600 hover:text-primary-700 p-1"
            title="Add User"
          />
          {role.name !== 'admin' && (
            <Button
              icon={<FaTrash className="w-3 h-3" />}
              onClick={() => onDelete(role.id)}
              variant="ghost"
              size="xs"
              className="text-primary-600 hover:text-primary-700 p-1"
              title="Delete Role"
            />
          )}
        </div>
      </div>
    </div>
  );
} 