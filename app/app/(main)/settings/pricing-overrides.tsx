import React, { useState } from 'react';
import { View, Text, ScrollView, TextInput, TouchableOpacity, Alert, ActivityIndicator } from 'react-native';
import { useRouter } from 'expo-router';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../../../src/api/client';

interface PricingOverride {
  id: string;
  override_type: string;
  item_key: string;
  override_value: number;
  is_percentage: boolean;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export default function PricingOverridesScreen() {
  const router = useRouter();
  const queryClient = useQueryClient();
  
  const [showAddForm, setShowAddForm] = useState(false);
  const [newOverride, setNewOverride] = useState({
    override_type: 'material',
    item_key: '',
    override_value: '',
    is_percentage: false,
    notes: '',
  });

  // Fetch pricing overrides
  const { data: overrides, isLoading, error } = useQuery<PricingOverride[]>({
    queryKey: ['pricingOverrides'],
    queryFn: async () => {
      const response = await api.get('/api/company/pricing-overrides');
      return response.data;
    },
  });

  // Create mutation
  const createMutation = useMutation({
    mutationFn: async (data: any) => {
      const response = await api.post('/api/company/pricing-overrides', {
        ...data,
        override_value: parseFloat(data.override_value),
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pricingOverrides'] });
      setShowAddForm(false);
      setNewOverride({
        override_type: 'material',
        item_key: '',
        override_value: '',
        is_percentage: false,
        notes: '',
      });
      Alert.alert('Success', 'Pricing override created successfully');
    },
    onError: (error: any) => {
      Alert.alert('Error', error.response?.data?.error || 'Failed to create override');
    },
  });

  // Delete mutation
  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/api/company/pricing-overrides/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pricingOverrides'] });
      Alert.alert('Success', 'Pricing override deleted');
    },
    onError: (error: any) => {
      Alert.alert('Error', error.response?.data?.error || 'Failed to delete override');
    },
  });

  const handleCreate = () => {
    if (!newOverride.item_key || !newOverride.override_value) {
      Alert.alert('Error', 'Please fill in all required fields');
      return;
    }
    createMutation.mutate(newOverride);
  };

  const handleDelete = (id: string) => {
    Alert.alert(
      'Confirm Delete',
      'Are you sure you want to delete this pricing override?',
      [
        { text: 'Cancel', style: 'cancel' },
        { text: 'Delete', style: 'destructive', onPress: () => deleteMutation.mutate(id) },
      ]
    );
  };

  if (isLoading) {
    return (
      <View className="flex-1 justify-center items-center bg-gray-50">
        <ActivityIndicator size="large" color="#3B82F6" />
      </View>
    );
  }

  if (error) {
    return (
      <View className="flex-1 justify-center items-center bg-gray-50 p-6">
        <Text className="text-red-600 text-center">Failed to load pricing overrides</Text>
      </View>
    );
  }

  return (
    <ScrollView className="flex-1 bg-gray-50">
      <View className="p-6">
        <View className="flex-row justify-between items-center mb-6">
          <View>
            <Text className="text-3xl font-bold text-gray-900">Pricing Overrides</Text>
            <Text className="text-gray-600 mt-1">
              Custom pricing for materials and labor
            </Text>
          </View>
        </View>

        {!showAddForm && (
          <TouchableOpacity
            onPress={() => setShowAddForm(true)}
            className="bg-blue-600 rounded-lg p-4 mb-6"
          >
            <Text className="text-white text-center font-semibold text-lg">
              + Add New Override
            </Text>
          </TouchableOpacity>
        )}

        {showAddForm && (
          <View className="bg-white rounded-lg p-6 mb-6 shadow-sm border border-gray-200">
            <Text className="text-xl font-bold text-gray-900 mb-4">New Override</Text>
            
            <Text className="text-sm font-medium text-gray-700 mb-2">Type</Text>
            <View className="flex-row mb-4">
              {['material', 'labor'].map((type) => (
                <TouchableOpacity
                  key={type}
                  onPress={() => setNewOverride({ ...newOverride, override_type: type })}
                  className={`flex-1 p-3 rounded-lg border mr-2 ${
                    newOverride.override_type === type
                      ? 'bg-blue-50 border-blue-600'
                      : 'bg-gray-50 border-gray-300'
                  }`}
                >
                  <Text
                    className={`text-center capitalize ${
                      newOverride.override_type === type ? 'text-blue-600 font-semibold' : 'text-gray-600'
                    }`}
                  >
                    {type}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>

            <Text className="text-sm font-medium text-gray-700 mb-2">
              Item Key (e.g., lumber, carpentry)
            </Text>
            <TextInput
              value={newOverride.item_key}
              onChangeText={(text) => setNewOverride({ ...newOverride, item_key: text })}
              placeholder="lumber"
              className="bg-gray-50 border border-gray-300 rounded-lg p-3 mb-4"
            />

            <Text className="text-sm font-medium text-gray-700 mb-2">
              Override Value
            </Text>
            <TextInput
              value={newOverride.override_value}
              onChangeText={(text) => setNewOverride({ ...newOverride, override_value: text })}
              placeholder="10.50"
              keyboardType="decimal-pad"
              className="bg-gray-50 border border-gray-300 rounded-lg p-3 mb-4"
            />

            <TouchableOpacity
              onPress={() => setNewOverride({ ...newOverride, is_percentage: !newOverride.is_percentage })}
              className="flex-row items-center mb-4"
            >
              <View className={`w-6 h-6 rounded border-2 mr-3 ${
                newOverride.is_percentage ? 'bg-blue-600 border-blue-600' : 'bg-white border-gray-300'
              }`}>
                {newOverride.is_percentage && (
                  <Text className="text-white text-center">✓</Text>
                )}
              </View>
              <Text className="text-gray-700">Is Percentage Adjustment</Text>
            </TouchableOpacity>

            <Text className="text-sm font-medium text-gray-700 mb-2">
              Notes (Optional)
            </Text>
            <TextInput
              value={newOverride.notes}
              onChangeText={(text) => setNewOverride({ ...newOverride, notes: text })}
              placeholder="Custom pricing notes"
              multiline
              numberOfLines={3}
              className="bg-gray-50 border border-gray-300 rounded-lg p-3 mb-4"
            />

            <View className="flex-row gap-3">
              <TouchableOpacity
                onPress={() => setShowAddForm(false)}
                className="flex-1 bg-gray-200 rounded-lg p-4"
              >
                <Text className="text-gray-700 text-center font-semibold">Cancel</Text>
              </TouchableOpacity>
              <TouchableOpacity
                onPress={handleCreate}
                disabled={createMutation.isPending}
                className={`flex-1 rounded-lg p-4 ${
                  createMutation.isPending ? 'bg-blue-400' : 'bg-blue-600'
                }`}
              >
                {createMutation.isPending ? (
                  <ActivityIndicator color="white" />
                ) : (
                  <Text className="text-white text-center font-semibold">Create</Text>
                )}
              </TouchableOpacity>
            </View>
          </View>
        )}

        {overrides && overrides.length === 0 ? (
          <View className="bg-white rounded-lg p-8 text-center">
            <Text className="text-gray-600 text-center mb-2">No pricing overrides yet</Text>
            <Text className="text-gray-500 text-center text-sm">
              Add custom pricing to override default material and labor costs
            </Text>
          </View>
        ) : (
          <View className="space-y-4">
            {overrides?.map((override) => (
              <View
                key={override.id}
                className="bg-white rounded-lg p-4 shadow-sm border border-gray-200"
              >
                <View className="flex-row justify-between items-start mb-3">
                  <View className="flex-1">
                    <View className="flex-row items-center mb-2">
                      <View className="bg-blue-100 px-3 py-1 rounded-full mr-2">
                        <Text className="text-blue-800 text-xs font-semibold uppercase">
                          {override.override_type}
                        </Text>
                      </View>
                      <Text className="text-lg font-bold text-gray-900">
                        {override.item_key}
                      </Text>
                    </View>
                    <Text className="text-2xl font-semibold text-green-600">
                      {override.is_percentage ? `${override.override_value}%` : `$${override.override_value.toFixed(2)}`}
                    </Text>
                    {override.notes && (
                      <Text className="text-sm text-gray-600 mt-2">{override.notes}</Text>
                    )}
                  </View>
                  <TouchableOpacity
                    onPress={() => handleDelete(override.id)}
                    disabled={deleteMutation.isPending}
                    className="ml-4"
                  >
                    <Text className="text-red-600 text-lg">🗑️</Text>
                  </TouchableOpacity>
                </View>
                <Text className="text-xs text-gray-400">
                  Updated {new Date(override.updated_at).toLocaleDateString()}
                </Text>
              </View>
            ))}
          </View>
        )}
      </View>
    </ScrollView>
  );
}
