import { useState } from 'react'

export default function Home() {
  const [count, setCount] = useState(0)

  return (
    <div>
      <h1>Home</h1>
      <button type="button" onClick={() => setCount((c) => c + 1)}>
        Count is {count}
      </button>
    </div>
  )
}
