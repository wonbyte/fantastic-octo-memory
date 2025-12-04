import React from 'react';
import { render } from '@testing-library/react-native';
import Index from '../app/index';

describe('Index screen', () => {
  it('renders correctly', () => {
    const { getByText } = render(<Index />);
    expect(getByText('Construction Estimator')).toBeTruthy();
    expect(getByText('AI-Powered Bidding Automation')).toBeTruthy();
  });
});
