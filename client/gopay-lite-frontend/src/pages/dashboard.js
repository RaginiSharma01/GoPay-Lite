import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import styles from '../styles/Dashboard.module.css';
import Navbar from '../components/Navbar';

export default function Dashboard() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [paying, setPaying] = useState(false);
  const [error, setError] = useState(null);
  const router = useRouter();

  // Fetch user details on mount
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login');
      return;
    }

    fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/auth/me`, {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then(async (res) => {
        if (!res.ok) {
          localStorage.removeItem('token');
          router.push('/login');
          return;
        }

        const data = await res.json();
        setUser({
          name: data.email || 'User',
          balance: data.balance || 1200,
          transactions: data.transactions || [],
          status: data.status || 'active'
        });
      })
      .catch((err) => {
        setError(err.message);
        localStorage.removeItem('token');
        router.push('/login');
      })
      .finally(() => setLoading(false));
  }, [router]);

  const handlePay = async () => {
    if (user?.balance < 500) {
      setError('Insufficient funds');
      return;
    }

    setPaying(true);
    setError(null);
    const token = localStorage.getItem('token');

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/pay`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          amount: 500,
          from_account: 'wallet',
          to_account: 'merchant',
        }),
      });

      const data = await res.json();

      if (res.ok) {
        alert(`Payment successful! Order ID: ${data.razorpay_order_id}`);
        // Update balance after successful payment
        setUser(prev => ({
          ...prev,
          balance: prev.balance - 500,
          transactions: [...prev.transactions, {
            id: data.razorpay_order_id,
            amount: -500,
            date: new Date().toISOString(),
            description: 'Payment to merchant'
          }]
        }));
      } else {
        setError(data.message || 'Payment failed');
      }
    } catch (err) {
      setError("Server error while processing payment.");
      console.error(err);
    } finally {
      setPaying(false);
    }
  };

  if (loading) return (
    <div className={styles.loadingContainer}>
      <Navbar />
      <div className={styles.spinner}></div>
    </div>
  );

  return (
    <>
      <Navbar />

      <div className={styles.container}>
        <div className={styles.header}>
          <h1>Welcome, {user?.name}</h1>
          {error && <div className={styles.error}>{error}</div>}
        </div>

        <div className={styles.cards}>
          <div className={styles.card}>
            <h3>Wallet Balance</h3>
            <p>₹{user?.balance}</p>
            {user?.balance < 500 && (
              <p className={styles.insufficientFunds}>Minimum ₹500 required for payment</p>
            )}
          </div>
          <div className={styles.card}>
            <h3>Transactions</h3>
            <p>{user?.transactions?.length || 0}</p>
          </div>
          <div className={styles.card}>
            <h3>Status</h3>
            <p className={user?.status === 'active' ? styles.activeStatus : ''}>
              {user?.status || 'Unknown'}
            </p>
          </div>
        </div>

        <div className={styles.transactions}>
          <h3>Recent Transactions</h3>
          <ul>
            {user?.transactions?.slice(0, 5).map(tx => (
              <li key={tx.id}>
                <span>{new Date(tx.date).toLocaleDateString()}</span>
                <span>{tx.amount > 0 ? '+' : ''}{tx.amount}</span>
              </li>
            ))}
            {user?.transactions?.length === 0 && <li>No transactions yet</li>}
          </ul>
        </div>

        <div style={{ textAlign: 'center', marginTop: '2rem' }}>
          <button
            className={styles.payButton}
            onClick={handlePay}
            disabled={paying || (user?.balance || 0) < 500}
          >
            {paying ? (
              <>
                <span className={styles.spinner}></span>
                Processing...
              </>
            ) : 'Pay ₹500'}
          </button>
        </div>

        <div className={styles.footerNote}>
          GoPay-Lite powered by Go + Microservices
        </div>
      </div>
    </>
  );
}