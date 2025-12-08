import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { useTheme, ThemeMode } from '../../contexts/ThemeContext';

export const ThemeToggle: React.FC = () => {
  const { themeMode, setThemeMode, colors } = useTheme();

  const modes: { value: ThemeMode; label: string; icon: string }[] = [
    { value: 'light', label: 'Light', icon: '☀️' },
    { value: 'dark', label: 'Dark', icon: '🌙' },
    { value: 'auto', label: 'Auto', icon: '⚙️' },
  ];

  return (
    <View style={styles.container}>
      <Text style={[styles.label, { color: colors.text }]}>Theme</Text>
      <View style={styles.buttonContainer}>
        {modes.map((mode) => (
          <TouchableOpacity
            key={mode.value}
            style={[
              styles.button,
              {
                backgroundColor:
                  themeMode === mode.value ? colors.primary : colors.surface,
                borderColor: colors.border,
              },
            ]}
            onPress={() => setThemeMode(mode.value)}
            accessibilityRole="button"
            accessibilityLabel={`${mode.label} theme`}
            accessibilityState={{ selected: themeMode === mode.value }}
          >
            <Text style={styles.icon}>{mode.icon}</Text>
            <Text
              style={[
                styles.buttonText,
                {
                  color:
                    themeMode === mode.value ? '#FFFFFF' : colors.text,
                },
              ]}
            >
              {mode.label}
            </Text>
          </TouchableOpacity>
        ))}
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginVertical: 16,
  },
  label: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 8,
  },
  buttonContainer: {
    flexDirection: 'row',
    gap: 8,
  },
  button: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 12,
    paddingHorizontal: 16,
    borderRadius: 8,
    borderWidth: 1,
    gap: 6,
  },
  icon: {
    fontSize: 20,
  },
  buttonText: {
    fontSize: 14,
    fontWeight: '600',
  },
});
