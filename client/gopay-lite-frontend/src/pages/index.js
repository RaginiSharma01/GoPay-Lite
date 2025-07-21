import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import styles from '../styles/Home.module.css';
import { motion, AnimatePresence } from 'framer-motion';
import Navbar from '../components/Navbar';

export default function Home() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const router = useRouter();

  // Step 1: Image list
  const illustrations = [
    '/forms.png',
    '/mobilepay.png',
    '/onlinepayments.png',
    '/payonline.png',
    '/profiledata.png'
  ];

  const [index, setIndex] = useState(0);

  // Step 2: Auto-rotate every 5 seconds
  useEffect(() => {
    const token = localStorage.getItem('token');
    setIsLoggedIn(!!token);

    const interval = setInterval(() => {
      setIndex((prev) => (prev + 1) % illustrations.length);
    }, 5000); // ⏱️ Change every 5 sec

    return () => clearInterval(interval);
  }, []);

  const goToDashboard = () => {
    router.push('/dashboard');
  };

  return (
    <>
      <Navbar />

      <div className={styles.hero}>
        <div className={styles.content}>
          {/* Animated image with switching */}
          <AnimatePresence mode="wait">
            <motion.img
              key={illustrations[index]}
              src={illustrations[index]}
              alt="Hero Illustration"
              className={styles.illustration}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.8 }}
            />
          </AnimatePresence>

          <motion.h1
            className={styles.title}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.3, duration: 0.8 }}
          >
            Welcome to <span className={styles.highlight}>GoPay-Lite</span>
          </motion.h1>

          <motion.p
            className={styles.subtitle}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.6, duration: 0.8 }}
          >
            Fast. Simple. Secure. Your modern microservice-powered payment gateway.
          </motion.p>

          <motion.div
            className={styles.buttonGroup}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1, duration: 0.8 }}
          >
            {!isLoggedIn ? (
              <>
                <Link href="/login">
                  <button className={styles.primaryButton}>Login</button>
                </Link>
                <Link href="/register">
                  <button className={styles.outlineButton}>Register</button>
                </Link>
              </>
            ) : (
              <button className={styles.primaryButton} onClick={goToDashboard}>
                Go to Dashboard
              </button>
            )}
          </motion.div>
        </div>
      </div>
    </>
  );
}
