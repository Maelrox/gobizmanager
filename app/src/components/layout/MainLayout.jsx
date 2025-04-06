import { useAuth } from '../../contexts/AuthContext';
import { route } from 'preact-router';
import { FaHome, FaBuilding, FaUsers } from 'react-icons/fa';

export default function MainLayout({ children }) {
  const { clearTokens } = useAuth();

  const handleLogout = () => {
    clearTokens();
    route('/login');
  };

  return (
    <div className="min-h-screen bg-neutral-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-2xl font-logo text-primary-600">GoBizManager</h1>
            </div>
            <div className="flex items-center">
              <button
                onClick={handleLogout}
                className="px-4 py-2 text-sm font-medium text-neutral-700 hover:text-primary-600 transition-colors duration-200"
              >
                Sign out
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex">
          {/* Sidebar */}
          <aside className="w-64 pr-4">
            <nav className="space-y-1">
              <a
                href="/dashboard"
                className="flex items-center px-4 py-2 text-sm font-medium text-neutral-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors duration-200"
              >
                <FaHome className="mr-3 h-5 w-5" />
                Dashboard
              </a>
              <a
                href="/companies"
                className="flex items-center px-4 py-2 text-sm font-medium text-neutral-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors duration-200"
              >
                <FaBuilding className="mr-3 h-5 w-5" />
                Companies
              </a>
              <a
                href="/roles"
                className="flex items-center px-4 py-2 text-sm font-medium text-neutral-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors duration-200"
              >
                <FaUsers className="mr-3 h-5 w-5" />
                Roles & Permissions
              </a>
            </nav>
          </aside>

          {/* Main Content Area */}
          <main className="flex-1">
            <div className="bg-white rounded-2xl shadow-sm p-6">
              {children}
            </div>
          </main>
        </div>
      </div>
    </div>
  );
} 