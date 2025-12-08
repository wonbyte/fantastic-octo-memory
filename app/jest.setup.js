/* eslint-env jest */
/* eslint-disable @typescript-eslint/no-require-imports */
/* eslint-disable no-undef */
import '@react-native-async-storage/async-storage/jest/async-storage-mock';

jest.mock('@react-native-async-storage/async-storage', () => 
  require('@react-native-async-storage/async-storage/jest/async-storage-mock')
);

jest.mock('@react-native-community/netinfo', () => ({
  addEventListener: jest.fn(() => jest.fn()),
  fetch: jest.fn(() => 
    Promise.resolve({
      isConnected: true,
      isInternetReachable: true,
    })
  ),
}));

jest.mock('expo-notifications', () => ({
  setNotificationHandler: jest.fn(),
  getPermissionsAsync: jest.fn(() => 
    Promise.resolve({ status: 'granted' })
  ),
  requestPermissionsAsync: jest.fn(() => 
    Promise.resolve({ status: 'granted' })
  ),
  getExpoPushTokenAsync: jest.fn(() => 
    Promise.resolve({ data: 'mock-token' })
  ),
  setNotificationChannelAsync: jest.fn(),
  scheduleNotificationAsync: jest.fn(() => 
    Promise.resolve('mock-notification-id')
  ),
  cancelAllScheduledNotificationsAsync: jest.fn(),
  addNotificationReceivedListener: jest.fn(() => ({
    remove: jest.fn(),
  })),
  addNotificationResponseReceivedListener: jest.fn(() => ({
    remove: jest.fn(),
  })),
  AndroidImportance: {
    MAX: 5,
  },
  SchedulableTriggerInputTypes: {
    TIME_INTERVAL: 'timeInterval',
  },
}));

jest.mock('expo-device', () => ({
  isDevice: true,
}));

jest.mock('expo-constants', () => ({
  expoConfig: {
    extra: {
      eas: {
        projectId: 'mock-project-id',
      },
    },
  },
}));
