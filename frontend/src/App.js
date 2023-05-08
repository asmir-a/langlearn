import React from 'react';
import { useState } from 'react';

const App = () => {
  const [count, setCount] = useState(0);
  return (
    <>
      <h1> Whatever {count} </h1>
      <button onClick = {() => setCount(count + 1)}> INC </button>
    </>
  );
}

export default App;
