import { Router, Route } from 'preact-router';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import MainLayout from './components/layout/MainLayout';
import LoginForm from './components/auth/LoginForm';
import RegisterForm from './components/auth/RegisterForm';
import CompanyList from './components/companies/CompanyList';
import PermissionsManager from './components/permissions/PermissionsManager';
import UserManager from './components/users/UserManager';

function PrivateRoute({ component: Component, ...rest }) {
  const { isAuthenticated } = useAuth();

  if (!isAuthenticated.value) {
    window.location.href = '/login';
    return null;
  }

  return (
    <MainLayout>
      <Component {...rest} />
    </MainLayout>
  );
}

function PublicRoute({ component: Component, ...rest }) {
  const { isAuthenticated } = useAuth();

  if (isAuthenticated.value) {
    window.location.href = '/dashboard';
    return null;
  }

  return <Component {...rest} />;
}

export function App() {
  return (
    <AuthProvider>
      <Router>
        <PublicRoute path="/login" component={LoginForm} />
        <PublicRoute path="/register" component={RegisterForm} />
        <PrivateRoute path="/companies" component={CompanyList} />
        <PrivateRoute path="/roles" component={PermissionsManager} />
        <PrivateRoute path="/users" component={UserManager} />
        <PrivateRoute path="/dashboard" component={CompanyList} />
        <PrivateRoute path="/" component={CompanyList} />
      </Router>
    </AuthProvider>
  );
} 