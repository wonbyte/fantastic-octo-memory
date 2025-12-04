import React from 'react';
import { render, fireEvent } from '@testing-library/react-native';
import { Button } from '../src/components/ui/Button';
import { Input } from '../src/components/ui/Input';
import { Card } from '../src/components/ui/Card';
import { Loading } from '../src/components/ui/Loading';
import { ErrorState } from '../src/components/ui/ErrorState';

describe('UI Components', () => {
  describe('Button', () => {
    it('renders with title', () => {
      const { getByText } = render(
        <Button title="Test Button" onPress={() => {}} />
      );
      expect(getByText('Test Button')).toBeTruthy();
    });

    it('calls onPress when pressed', () => {
      const onPress = jest.fn();
      const { getByText } = render(
        <Button title="Test Button" onPress={onPress} />
      );
      fireEvent.press(getByText('Test Button'));
      expect(onPress).toHaveBeenCalledTimes(1);
    });

    it('does not call onPress when disabled', () => {
      const onPress = jest.fn();
      const { getByText } = render(
        <Button title="Test Button" onPress={onPress} disabled />
      );
      fireEvent.press(getByText('Test Button'));
      expect(onPress).not.toHaveBeenCalled();
    });

    it('shows loading indicator when loading', () => {
      const { queryByText, UNSAFE_root } = render(
        <Button title="Test Button" onPress={() => {}} loading />
      );
      expect(queryByText('Test Button')).toBeFalsy();
      expect(UNSAFE_root.findAllByType('ActivityIndicator')).toHaveLength(1);
    });
  });

  describe('Input', () => {
    it('renders with label', () => {
      const { getByText } = render(
        <Input label="Test Label" value="" onChangeText={() => {}} />
      );
      expect(getByText('Test Label')).toBeTruthy();
    });

    it('displays error message', () => {
      const { getByText } = render(
        <Input
          label="Test Label"
          value=""
          onChangeText={() => {}}
          error="This is an error"
        />
      );
      expect(getByText('This is an error')).toBeTruthy();
    });

    it('calls onChangeText when text changes', () => {
      const onChangeText = jest.fn();
      const { getByDisplayValue } = render(
        <Input
          label="Test Label"
          value="test"
          onChangeText={onChangeText}
        />
      );
      const input = getByDisplayValue('test');
      fireEvent.changeText(input, 'new text');
      expect(onChangeText).toHaveBeenCalledWith('new text');
    });
  });

  describe('Card', () => {
    it('renders children', () => {
      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const { Text } = require('react-native');
      const { getByText } = render(
        <Card>
          <Text>Card Content</Text>
        </Card>
      );
      expect(getByText('Card Content')).toBeTruthy();
    });
  });

  describe('Loading', () => {
    it('renders with default message', () => {
      const { getByText } = render(<Loading />);
      expect(getByText('Loading...')).toBeTruthy();
    });

    it('renders with custom message', () => {
      const { getByText } = render(<Loading message="Custom loading" />);
      expect(getByText('Custom loading')).toBeTruthy();
    });
  });

  describe('ErrorState', () => {
    it('renders with default message', () => {
      const { getByText } = render(<ErrorState />);
      expect(getByText('Something went wrong')).toBeTruthy();
    });

    it('renders with custom message', () => {
      const { getByText } = render(
        <ErrorState message="Custom error message" />
      );
      expect(getByText('Custom error message')).toBeTruthy();
    });

    it('shows retry button when onRetry is provided', () => {
      const onRetry = jest.fn();
      const { getByText } = render(<ErrorState onRetry={onRetry} />);
      expect(getByText('Try Again')).toBeTruthy();
    });

    it('calls onRetry when retry button is pressed', () => {
      const onRetry = jest.fn();
      const { getByText } = render(<ErrorState onRetry={onRetry} />);
      fireEvent.press(getByText('Try Again'));
      expect(onRetry).toHaveBeenCalledTimes(1);
    });
  });
});
