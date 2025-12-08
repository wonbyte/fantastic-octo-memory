import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Button } from './Button';
import { useTheme } from '../../contexts/ThemeContext';

interface ErrorStateProps {
  message?: string;
  onRetry?: () => void;
}

export const ErrorState: React.FC<ErrorStateProps> = ({
  message = 'Something went wrong',
  onRetry,
}) => {
  const { colors } = useTheme();
  
  return (
    <View
      style={[styles.container, { backgroundColor: colors.background }]}
      accessibilityRole="alert"
    >
      <Text style={styles.emoji} accessibilityLabel="Error icon">
        ⚠️
      </Text>
      <Text style={[styles.message, { color: colors.textSecondary }]}>
        {message}
      </Text>
      {onRetry && (
        <Button
          title="Try Again"
          onPress={onRetry}
          style={styles.button}
          accessibilityHint="Retry the failed operation"
        />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
  },
  emoji: {
    fontSize: 48,
    marginBottom: 16,
  },
  message: {
    fontSize: 16,
    textAlign: 'center',
    marginBottom: 24,
  },
  button: {
    minWidth: 120,
  },
});
