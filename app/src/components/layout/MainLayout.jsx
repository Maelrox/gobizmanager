import { useAuth } from '../../contexts/AuthContext';
import { route } from 'preact-router';
import { FaHome, FaBuilding, FaUsers, FaBars, FaTimes, FaSignOutAlt, FaUser } from 'react-icons/fa';
import { useState, useEffect } from 'preact/hooks';

export default function MainLayout({ children }) {
  const { clearTokens } = useAuth();
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
      if (window.innerWidth >= 768) {
        setIsSidebarOpen(true);
      } else {
        setIsSidebarOpen(false);
      }
    };

    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  const handleLogout = () => {
    clearTokens();
    route('/login');
  };

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  return (
    <div className="min-h-screen bg-neutral-50 flex flex-col">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-full mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-2xl font-logo text-primary-600">GoBizManager</h1>
            </div>
            <div className="flex items-center">
              <button
                onClick={toggleSidebar}
                className="p-2 rounded-lg hover:bg-neutral-100 transition-colors duration-200 md:hidden"
              >
                <FaBars className="h-5 w-5 text-neutral-600" />
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar */}
        <aside
          className={`fixed md:relative inset-y-0 left-0 z-50 bg-white shadow-lg transition-all duration-300 ease-in-out ${
            isSidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'
          } ${isMobile ? 'w-full md:w-64' : 'w-64'}`}
        >
          <div className="flex justify-between items-center p-4 border-b border-neutral-200">
            <h2 className="text-lg font-semibold text-neutral-900">Menu</h2>
            <button
              onClick={toggleSidebar}
              className="p-2 rounded-lg hover:bg-neutral-100 transition-colors duration-200 md:hidden"
            >
              <FaTimes className="h-5 w-5 text-neutral-600" />
            </button>
          </div>
          <nav className="h-full py-4">
            <div className="space-y-1 px-4">
              <a
                href="/dashboard"
                onClick={() => isMobile && setIsSidebarOpen(false)}
                className="flex items-center px-4 py-2 text-sm font-medium text-neutral-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors duration-200"
              >
                <FaHome className="mr-3 h-5 w-5" />
                <span>Dashboard</span>
              </a>
              <a
                href="/companies"
                onClick={() => isMobile && setIsSidebarOpen(false)}
                className="flex items-center px-4 py-2 text-sm font-medium text-neutral-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors duration-200"
              >
                <FaBuilding className="mr-3 h-5 w-5" />
                <span>Companies</span>
              </a>
              <a
                href="/roles"
                onClick={() => isMobile && setIsSidebarOpen(false)}
                className="flex items-center px-4 py-2 text-sm font-medium text-neutral-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors duration-200"
              >
                <FaUsers className="mr-3 h-5 w-5" />
                <span>Roles & Permissions</span>
              </a>
              <a
                href="/users"
                onClick={() => isMobile && setIsSidebarOpen(false)}
                className="flex items-center px-4 py-2 text-sm font-medium text-neutral-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors duration-200"
              >
                <FaUser className="mr-3 h-5 w-5" />
                <span>Users</span>
              </a>
              <button
                onClick={() => {
                  if (isMobile) {
                    setIsSidebarOpen(false);
                  }
                  handleLogout();
                }}
                className="w-full flex items-center px-4 py-2 text-sm font-medium text-neutral-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors duration-200"
              >
                <FaSignOutAlt className="mr-3 h-5 w-5" />
                <span>Sign out</span>
              </button>
            </div>
          </nav>
        </aside>

        {/* Overlay for mobile */}
        {isMobile && isSidebarOpen && (
          <div
            className="fixed inset-0 bg-black bg-opacity-50 z-40 md:hidden"
            onClick={toggleSidebar}
          />
        )}

        {/* Main Content Area */}
        <main className="flex-1 overflow-auto">
          <div className="p-6">
            <div className="bg-white rounded-2xl shadow-sm p-6">
              {children}
            </div>
          </div>
        </main>
      </div>
    </div>
  );
} 