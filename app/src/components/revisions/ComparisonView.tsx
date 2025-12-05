import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  ActivityIndicator,
} from 'react-native';
import { Card } from '../ui/Card';
import { COLORS } from '../../utils/constants';
import {
  BlueprintComparison,
  BidComparison,
  BlueprintChange,
  BidChange,
} from '../../types';
import { revisionsApi } from '../../api/revisions';

type ComparisonType = 'blueprint' | 'bid';

interface ComparisonViewProps {
  itemId: string;
  type: ComparisonType;
  fromVersion: number;
  toVersion: number;
}

export const ComparisonView: React.FC<ComparisonViewProps> = ({
  itemId,
  type,
  fromVersion,
  toVersion,
}) => {
  const [comparison, setComparison] = useState<BlueprintComparison | BidComparison | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadComparison();
  }, [itemId, type, fromVersion, toVersion]);

  const loadComparison = async () => {
    try {
      setLoading(true);
      setError(null);
      
      let data: BlueprintComparison | BidComparison;
      if (type === 'blueprint') {
        data = await revisionsApi.compareBlueprintRevisions(itemId, fromVersion, toVersion);
      } else {
        data = await revisionsApi.compareBidRevisions(itemId, fromVersion, toVersion);
      }
      
      setComparison(data);
    } catch (err) {
      setError('Failed to load comparison');
      console.error('Error loading comparison:', err);
    } finally {
      setLoading(false);
    }
  };

  const getChangeIcon = (changeType: string) => {
    switch (changeType) {
      case 'added':
        return '+';
      case 'removed':
        return '-';
      case 'modified':
        return '~';
      default:
        return '?';
    }
  };

  const getChangeColor = (changeType: string) => {
    switch (changeType) {
      case 'added':
        return COLORS.success;
      case 'removed':
        return COLORS.error;
      case 'modified':
        return COLORS.warning;
      default:
        return COLORS.text.secondary;
    }
  };

  const getImpactColor = (impact?: string) => {
    switch (impact) {
      case 'High':
        return COLORS.error;
      case 'Medium':
        return COLORS.warning;
      case 'Low':
        return COLORS.primary;
      default:
        return COLORS.text.secondary;
    }
  };

  const renderChange = (change: BlueprintChange | BidChange, index: number) => {
    const isBidChange = 'trade' in change;
    
    return (
      <Card key={index} style={styles.changeCard}>
        <View style={styles.changeHeader}>
          <View style={styles.changeTypeContainer}>
            <View
              style={[
                styles.changeIcon,
                { backgroundColor: getChangeColor(change.change_type) },
              ]}
            >
              <Text style={styles.changeIconText}>{getChangeIcon(change.change_type)}</Text>
            </View>
            <Text style={styles.changeCategory}>{change.category}</Text>
          </View>
          
          {change.impact && (
            <View
              style={[
                styles.impactBadge,
                { backgroundColor: getImpactColor(change.impact) },
              ]}
            >
              <Text style={styles.impactText}>{change.impact}</Text>
            </View>
          )}
        </View>
        
        <Text style={styles.changeDescription}>{change.description}</Text>
        
        {isBidChange && (change as BidChange).trade && (
          <Text style={styles.tradeBadgeText}>Trade: {(change as BidChange).trade}</Text>
        )}
      </Card>
    );
  };

  if (loading) {
    return (
      <View style={styles.centerContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
        <Text style={styles.loadingText}>Analyzing changes...</Text>
      </View>
    );
  }

  if (error || !comparison) {
    return (
      <View style={styles.centerContainer}>
        <Text style={styles.errorText}>{error || 'No comparison data available'}</Text>
      </View>
    );
  }

  return (
    <ScrollView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>
          Comparing Version {fromVersion} → {toVersion}
        </Text>
      </View>

      <Card style={styles.summaryCard}>
        <Text style={styles.summaryTitle}>Summary</Text>
        <View style={styles.summaryRow}>
          <Text style={styles.summaryLabel}>Total Changes:</Text>
          <Text style={styles.summaryValue}>{comparison.summary.total_changes}</Text>
        </View>
        <View style={styles.summaryRow}>
          <Text style={[styles.summaryLabel, { color: COLORS.success }]}>Added:</Text>
          <Text style={styles.summaryValue}>{comparison.summary.added_count}</Text>
        </View>
        <View style={styles.summaryRow}>
          <Text style={[styles.summaryLabel, { color: COLORS.warning }]}>Modified:</Text>
          <Text style={styles.summaryValue}>{comparison.summary.modified_count}</Text>
        </View>
        <View style={styles.summaryRow}>
          <Text style={[styles.summaryLabel, { color: COLORS.error }]}>Removed:</Text>
          <Text style={styles.summaryValue}>{comparison.summary.removed_count}</Text>
        </View>
        {comparison.summary.high_impact_count > 0 && (
          <View style={styles.summaryRow}>
            <Text style={[styles.summaryLabel, { color: COLORS.error }]}>High Impact:</Text>
            <Text style={[styles.summaryValue, { fontWeight: 'bold' }]}>
              {comparison.summary.high_impact_count}
            </Text>
          </View>
        )}
      </Card>

      <View style={styles.changesSection}>
        <Text style={styles.sectionTitle}>Changes</Text>
        {comparison.changes.length === 0 ? (
          <Card>
            <Text style={styles.noChangesText}>No changes detected</Text>
          </Card>
        ) : (
          comparison.changes.map(renderChange)
        )}
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.secondary,
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  header: {
    padding: 16,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.text.primary,
  },
  loadingText: {
    marginTop: 12,
    fontSize: 14,
    color: COLORS.text.secondary,
  },
  errorText: {
    fontSize: 14,
    color: COLORS.error,
    textAlign: 'center',
  },
  summaryCard: {
    margin: 16,
    marginTop: 0,
  },
  summaryTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: COLORS.text.primary,
    marginBottom: 12,
  },
  summaryRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingVertical: 4,
  },
  summaryLabel: {
    fontSize: 14,
    color: COLORS.text.secondary,
  },
  summaryValue: {
    fontSize: 14,
    color: COLORS.text.primary,
    fontWeight: '500',
  },
  changesSection: {
    padding: 16,
    paddingTop: 0,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: COLORS.text.primary,
    marginBottom: 12,
  },
  changeCard: {
    marginBottom: 12,
  },
  changeHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  changeTypeContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  changeIcon: {
    width: 24,
    height: 24,
    borderRadius: 12,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 8,
  },
  changeIconText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
  changeCategory: {
    fontSize: 14,
    fontWeight: '600',
    color: COLORS.text.primary,
    textTransform: 'capitalize',
  },
  impactBadge: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
  },
  impactText: {
    fontSize: 12,
    color: '#fff',
    fontWeight: '500',
  },
  changeDescription: {
    fontSize: 14,
    color: COLORS.text.secondary,
    lineHeight: 20,
  },
  tradeBadgeText: {
    fontSize: 12,
    color: COLORS.text.secondary,
    marginTop: 8,
    fontStyle: 'italic',
  },
  noChangesText: {
    fontSize: 14,
    color: COLORS.text.secondary,
    textAlign: 'center',
    padding: 20,
  },
});
