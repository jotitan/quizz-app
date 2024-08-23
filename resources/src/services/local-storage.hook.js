import { useEffect, useState } from 'react'

function useLocalStorage(localStorageKey, initialState){
  let localItem = localStorage.getItem(localStorageKey);
  let previousState = localItem ? JSON.parse(localItem) : initialState;
  const [value, setValue] = useState(previousState);

  useEffect(() => {
    if (value) {
      localStorage.setItem(localStorageKey, JSON.stringify(value))
    } else {
      localStorage.removeItem(localStorageKey)
    }
  }, [localStorageKey, value]);

  return [value, setValue]
}

export default useLocalStorage
