import Link from 'next/link';
import styles from '../styles/Navbar.module.css';
import { useContext, useEffect, useState } from 'react';
import { ThemeContext } from '../pages/_app';
import { useRouter } from 'next/router';

export default function Navbar() {
  const { darkMode, toggleTheme } = useContext(ThemeContext);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem('token');
    setIsLoggedIn(!!token); // Set based on presence of token
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    setIsLoggedIn(false);
    router.push('/login');
  };

  return (
    <nav className={styles.navbar}>
      <div className={styles.logo}>
        <Link href="/">GoPay-Lite</Link>
      </div>
      <div className={styles.links}>
        <Link href="/">Home</Link>
        {!isLoggedIn && (
          <>
            <Link href="/login">Login</Link>
            <Link href="/register">Register</Link>
          </>
        )}
        {isLoggedIn && (
          <>
            <Link href="/dashboard">Dashboard</Link>
            <button onClick={handleLogout} className={styles.logoutBtn}>
              Logout
            </button>
          </>
        )}
        <button className={styles.themeToggle} onClick={toggleTheme}>
          {darkMode ? '☀️' : '🌙'}
        </button>
      </div>
    </nav>
  );
}
