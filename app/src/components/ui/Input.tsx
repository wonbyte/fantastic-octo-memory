import React from 'react';
import {
  TextInput,
  Text,
  View,
  StyleSheet,
  TextInputProps,
  ViewStyle,
} from 'react-native';
import { useTheme } from '../../contexts/ThemeContext';

interface InputProps extends TextInputProps {
  label?: string;
  error?: string;
  containerStyle?: ViewStyle;
}

export const Input: React.FC<InputProps> = ({
  label,
  error,
  containerStyle,
  style,
  ...props
}) => {
  const { colors } = useTheme();
  
  return (
    <View style={[styles.container, containerStyle]}>
      {label && (
        <Text
          style={[styles.label, { color: colors.text }]}
          accessibilityLabel={label}
        >
          {label}
        </Text>
      )}
      <TextInput
        style={[
          styles.input,
          {
            borderColor: error ? colors.error : colors.border,
            color: colors.text,
            backgroundColor: colors.background,
          },
          style,
        ]}
        placeholderTextColor={colors.textLight}
        accessibilityLabel={label || 'Text input'}
        accessibilityHint={error}
        {...props}
      />
      {error && (
        <Text
          style={[styles.errorText, { color: colors.error }]}
          accessibilityRole="alert"
        >
          {error}
        </Text>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginBottom: 16,
  },
  label: {
    fontSize: 14,
    fontWeight: '600',
    marginBottom: 8,
  },
  input: {
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 16,
    paddingVertical: 12,
    fontSize: 16,
    minHeight: 48,
  },
  errorText: {
    fontSize: 12,
    marginTop: 4,
  },
});
