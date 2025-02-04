import { useState, useEffect} from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [items, setItems] = useState([]);

  useEffect(() => {
      fetch("http://localhost:8080/api/items")
          .then((res) => res.json())
          .then((data) => setItems(data))
          .catch((err) => console.error("Error fetching data:", err));
  }, []);

  return (
      <div>
          <h1>OSRS Flipping Assistant</h1>
          <ul>
              {items.map((item) => (
                  <li key={item.id}>
                  {item.name} (ID: {item.id}), Buy: {item.buy}, Sell: {item.sell}, Margin: {item.margin}
                  </li>
              ))}
          </ul>
      </div>
  );
}

export default App;
