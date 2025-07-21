export const makePayment = async (amount) => {
  const token = localStorage.getItem('token');

  if (!token) {
    throw new Error('User not authenticated. Token missing.');
  }

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

    if (!res.ok) {
      throw new Error(data.message || 'Payment failed');
    }

    return data;
  } catch (err) {
    console.error('Payment error:', err.message);
    throw err;
  }
};
