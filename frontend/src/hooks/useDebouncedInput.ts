import { useState, useCallback, useRef, useEffect } from 'react';

interface UseDebouncedInputOptions {
  delay?: number;
  minLength?: number;
  onDebouncedChange?: (value: string) => void;
}

interface DebouncedInputState {
  value: string;
  debouncedValue: string;
  isDebouncing: boolean;
}

interface DebouncedInputActions {
  setValue: (value: string) => void;
  clearValue: () => void;
  forceUpdate: () => void;
}

export const useDebouncedInput = ({
  delay = 300,
  minLength = 0,
  onDebouncedChange,
}: UseDebouncedInputOptions = {}): [DebouncedInputState, DebouncedInputActions] => {
  const [value, setValue] = useState('');
  const [debouncedValue, setDebouncedValue] = useState('');
  const [isDebouncing, setIsDebouncing] = useState(false);
  
  const timeoutRef = useRef<NodeJS.Timeout | undefined>(undefined);

  const updateDebouncedValue = useCallback((newValue: string) => {
    setDebouncedValue(newValue);
    setIsDebouncing(false);
    
    if (onDebouncedChange && newValue.length >= minLength) {
      onDebouncedChange(newValue);
    }
  }, [onDebouncedChange, minLength]);

  const handleSetValue = useCallback((newValue: string) => {
    setValue(newValue);
    
    // Clear existing timeout
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    // Don't debounce if value is empty or below minimum length
    if (newValue.length === 0 || newValue.length < minLength) {
      updateDebouncedValue(newValue);
      return;
    }

    setIsDebouncing(true);
    
    // Set new timeout
    timeoutRef.current = setTimeout(() => {
      updateDebouncedValue(newValue);
    }, delay);
  }, [delay, minLength, updateDebouncedValue]);

  const clearValue = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    setValue('');
    setDebouncedValue('');
    setIsDebouncing(false);
  }, []);

  const forceUpdate = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    updateDebouncedValue(value);
  }, [value, updateDebouncedValue]);

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return [
    {
      value,
      debouncedValue,
      isDebouncing,
    },
    {
      setValue: handleSetValue,
      clearValue,
      forceUpdate,
    },
  ];
};

export default useDebouncedInput;