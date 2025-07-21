import '../styles/globals.css';
import { useState, createContext } from 'react';

export const ThemeContext = createContext();

function MyApp({ Component, pageProps }) {
  const [darkMode, setDarkMode] = useState(false);

  const toggleTheme = () => {
    setDarkMode(!darkMode);
    document.documentElement.setAttribute(
      'data-theme',
      !darkMode ? 'dark' : 'light'
    );
  };

  return (
    <ThemeContext.Provider value={{ darkMode, toggleTheme }}>
      <Component {...pageProps} />
    </ThemeContext.Provider>
  );
}

export default MyApp;
