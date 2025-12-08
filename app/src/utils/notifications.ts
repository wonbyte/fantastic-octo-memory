import * as Notifications from 'expo-notifications';
import * as Device from 'expo-device';
import { Platform } from 'react-native';
import Constants from 'expo-constants';

// Configure notification handler
Notifications.setNotificationHandler({
  handleNotification: async () => ({
    shouldShowAlert: true,
    shouldPlaySound: true,
    shouldSetBadge: true,
    shouldShowBanner: true,
    shouldShowList: true,
  }),
});

export interface NotificationData {
  title: string;
  body: string;
  data?: Record<string, unknown>;
}

/**
 * Request notification permissions and get push token
 */
export async function registerForPushNotificationsAsync(): Promise<string | undefined> {
  if (!Device.isDevice) {
    console.log('Push notifications only work on physical devices');
    return undefined;
  }

  try {
    const { status: existingStatus } = await Notifications.getPermissionsAsync();
    let finalStatus = existingStatus;

    if (existingStatus !== 'granted') {
      const { status } = await Notifications.requestPermissionsAsync();
      finalStatus = status;
    }

    if (finalStatus !== 'granted') {
      console.log('Failed to get push token for push notification!');
      return undefined;
    }

    const projectId = Constants.expoConfig?.extra?.eas?.projectId;
    const token = await Notifications.getExpoPushTokenAsync({
      projectId: projectId || undefined,
    });

    if (Platform.OS === 'android') {
      await Notifications.setNotificationChannelAsync('default', {
        name: 'default',
        importance: Notifications.AndroidImportance.MAX,
        vibrationPattern: [0, 250, 250, 250],
        lightColor: '#3B82F6',
      });
    }

    return token.data;
  } catch (error) {
    console.error('Error registering for push notifications:', error);
    return undefined;
  }
}

/**
 * Schedule a local notification
 * @param notification - Notification content
 * @param delaySeconds - Delay in seconds before showing notification (0 = immediate)
 */
export async function scheduleLocalNotification(
  notification: NotificationData,
  delaySeconds: number = 0
): Promise<string> {
  const trigger = delaySeconds > 0 
    ? { 
        type: Notifications.SchedulableTriggerInputTypes.TIME_INTERVAL as const,
        seconds: delaySeconds,
        repeats: false,
      }
    : null;

  return await Notifications.scheduleNotificationAsync({
    content: {
      title: notification.title,
      body: notification.body,
      data: notification.data,
      sound: true,
    },
    trigger,
  });
}

/**
 * Send notification for job completion
 */
export async function notifyJobComplete(jobType: string, success: boolean): Promise<void> {
  await scheduleLocalNotification({
    title: success ? '✅ Job Complete' : '❌ Job Failed',
    body: success
      ? `Your ${jobType} has completed successfully`
      : `Your ${jobType} has failed. Please try again.`,
    data: { jobType, success },
  });
}

/**
 * Send notification for blueprint analysis
 */
export async function notifyAnalysisComplete(blueprintName: string): Promise<void> {
  await scheduleLocalNotification({
    title: '🎉 Analysis Complete',
    body: `Analysis for "${blueprintName}" is ready to view`,
    data: { type: 'analysis_complete', blueprintName },
  });
}

/**
 * Send notification for error
 */
export async function notifyError(message: string): Promise<void> {
  await scheduleLocalNotification({
    title: '⚠️ Error',
    body: message,
    data: { type: 'error' },
  });
}

/**
 * Cancel all scheduled notifications
 */
export async function cancelAllNotifications(): Promise<void> {
  await Notifications.cancelAllScheduledNotificationsAsync();
}

/**
 * Get notification permission status
 */
export async function getNotificationPermissionStatus(): Promise<string> {
  const { status } = await Notifications.getPermissionsAsync();
  return status;
}
