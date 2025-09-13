import { configureStore } from '@reduxjs/toolkit';
import chatReducer from './slices/chatSlice';
import connectionReducer from './slices/connectionSlice';

export const store = configureStore({
  reducer: {
    chat: chatReducer,
    connection: connectionReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // All state is now serializable with polling-based chat
      },
    }),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;