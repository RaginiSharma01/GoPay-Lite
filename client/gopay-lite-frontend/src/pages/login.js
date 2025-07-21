import { useState } from 'react';
import { useRouter } from 'next/router';
import { login } from '../services/auth';
import styles from '../styles/login.module.css';

export default function Login() {
  const router = useRouter();
  const [form, setForm] = useState({ email: '', password: '' });
  const [message, setMessage] = useState({ text: '', type: '' }); // success/error
  const [isLoading, setIsLoading] = useState(false);

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage({ text: '', type: '' });
    setIsLoading(true);

    try {
      const { token } = await login(form);
      
      // Secure token storage (consider httpOnly cookies for production)
      localStorage.setItem('token', token);
      
      // Redirect with success state
      await router.push('/dashboard?login=success');
    } catch (err) {
      setMessage({
        text: err.message || 'Login failed. Please try again.',
        type: 'error'
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={styles.loginContainer}>
      <h1 className={styles.heading}>Welcome Back</h1>
      <form onSubmit={handleSubmit} className={styles.form}>
        <div className={styles.inputGroup}>
          <label htmlFor="email" className={styles.label}>
            Email
          </label>
          <input
            id="email"
            name="email"
            type="email"
            autoComplete="username"
            placeholder="your@email.com"
            value={form.email}
            onChange={handleChange}
            required
            className={styles.input}
          />
        </div>

        <div className={styles.inputGroup}>
          <label htmlFor="password" className={styles.label}>
            Password
          </label>
          <input
            id="password"
            name="password"
            type="password"
            autoComplete="current-password"
            placeholder="••••••••"
            value={form.password}
            onChange={handleChange}
            required
            minLength={8}
            className={styles.input}
          />
        </div>

        <button
          type="submit"
          disabled={isLoading}
          className={`${styles.primaryButton} ${isLoading ? styles.loading : ''}`}
        >
          {isLoading ? 'Signing in...' : 'Sign In'}
        </button>

        {message.text && (
          <p className={`${styles.message} ${styles[message.type]}`}>
            {message.text}
          </p>
        )}

        <div className={styles.secondaryActions}>
          <a href="/forgot-password" className={styles.link}>
            Forgot password?
          </a>
          <a href="/register" className={styles.link}>
            Create account
          </a>
        </div>
      </form>
    </div>
  );
}