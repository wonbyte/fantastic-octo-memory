import React from 'react';
import { View, Text, ScrollView, TouchableOpacity } from 'react-native';
import { useRouter } from 'expo-router';
import type { Href } from 'expo-router';

export default function SettingsScreen() {
  const router = useRouter();

  const settingsOptions = [
    {
      title: 'Pricing Overrides',
      description: 'Manage custom material and labor pricing',
      route: '/settings/pricing-overrides' as Href,
      icon: '💰',
    },
    {
      title: 'Regional Settings',
      description: 'View regional pricing adjustments',
      route: '/settings/regional-adjustments' as Href,
      icon: '🌍',
    },
    {
      title: 'Cost Database Sync',
      description: 'Sync pricing data from external providers',
      route: '/settings/sync-costs' as Href,
      icon: '🔄',
    },
  ];

  return (
    <ScrollView className="flex-1 bg-gray-50">
      <View className="p-6">
        <Text className="text-3xl font-bold text-gray-900 mb-2">Settings</Text>
        <Text className="text-gray-600 mb-6">
          Manage your pricing configuration and cost database
        </Text>

        <View className="space-y-4">
          {settingsOptions.map((option, index) => (
            <TouchableOpacity
              key={index}
              onPress={() => router.push(option.route)}
              className="bg-white rounded-lg p-4 shadow-sm border border-gray-200"
            >
              <View className="flex-row items-center">
                <Text className="text-3xl mr-4">{option.icon}</Text>
                <View className="flex-1">
                  <Text className="text-lg font-semibold text-gray-900">
                    {option.title}
                  </Text>
                  <Text className="text-sm text-gray-600 mt-1">
                    {option.description}
                  </Text>
                </View>
                <Text className="text-gray-400 text-xl">›</Text>
              </View>
            </TouchableOpacity>
          ))}
        </View>

        <View className="mt-8 p-4 bg-blue-50 rounded-lg border border-blue-200">
          <Text className="text-sm text-blue-900 font-medium mb-2">
            ℹ️ About Cost Database Integration
          </Text>
          <Text className="text-sm text-blue-800">
            The cost database integrates with RSMeans, Home Depot, and Lowes to
            provide real-time pricing. Custom overrides allow you to adjust
            pricing specific to your business needs.
          </Text>
        </View>
      </View>
    </ScrollView>
  );
}
