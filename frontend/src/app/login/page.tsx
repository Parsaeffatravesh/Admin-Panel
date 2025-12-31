import { LoginProviders } from '@/components/LoginProviders';
import LoginForm from './LoginForm';

export default function LoginPage() {
  return (
    <LoginProviders>
      <LoginForm />
    </LoginProviders>
  );
}
