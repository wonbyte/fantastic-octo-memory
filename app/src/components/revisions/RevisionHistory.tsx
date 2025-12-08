import React, { useState, useEffect, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native';
import { Card } from '../ui/Card';
import { COLORS } from '../../utils/constants';
import { BlueprintRevision, BidRevision } from '../../types';
import { revisionsApi } from '../../api/revisions';

type RevisionType = 'blueprint' | 'bid';

interface RevisionHistoryProps {
  itemId: string;
  type: RevisionType;
  onCompare?: (fromVersion: number, toVersion: number) => void;
}

export const RevisionHistory: React.FC<RevisionHistoryProps> = ({
  itemId,
  type,
  onCompare,
}) => {
  const [revisions, setRevisions] = useState<(BlueprintRevision | BidRevision)[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedVersions, setSelectedVersions] = useState<number[]>([]);

  const loadRevisions = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      let data: (BlueprintRevision | BidRevision)[];
      if (type === 'blueprint') {
        data = await revisionsApi.getBlueprintRevisions(itemId);
      } else {
        data = await revisionsApi.getBidRevisions(itemId);
      }
      
      setRevisions(data);
    } catch (err) {
      setError('Failed to load revision history');
      console.error('Error loading revisions:', err);
    } finally {
      setLoading(false);
    }
  }, [itemId, type]);

  useEffect(() => {
    loadRevisions();
  }, [loadRevisions]);

  const handleVersionSelect = (version: number) => {
    if (selectedVersions.includes(version)) {
      setSelectedVersions(selectedVersions.filter((v) => v !== version));
    } else if (selectedVersions.length < 2) {
      setSelectedVersions([...selectedVersions, version]);
    }
  };

  const handleCompare = () => {
    if (selectedVersions.length === 2 && onCompare) {
      const [from, to] = selectedVersions.sort((a, b) => a - b);
      onCompare(from, to);
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const renderRevision = ({ item }: { item: BlueprintRevision | BidRevision }) => {
    const isSelected = selectedVersions.includes(item.version);
    const isBidRevision = 'final_price' in item;

    return (
      <TouchableOpacity
        onPress={() => handleVersionSelect(item.version)}
        disabled={!onCompare}
      >
        <Card style={isSelected ? StyleSheet.flatten([styles.revisionCard, styles.selectedCard]) : styles.revisionCard}>
          <View style={styles.revisionHeader}>
            <Text style={styles.versionText}>Version {item.version}</Text>
            <Text style={styles.dateText}>{formatDate(item.created_at)}</Text>
          </View>
          
          {isBidRevision && (
            <View style={styles.revisionDetails}>
              <Text style={styles.detailText}>
                Price: ${((item as BidRevision).final_price || 0).toLocaleString()}
              </Text>
            </View>
          )}
          
          {!isBidRevision && (
            <View style={styles.revisionDetails}>
              <Text style={styles.detailText}>
                File: {(item as BlueprintRevision).filename}
              </Text>
            </View>
          )}
          
          {item.changes_summary && (
            <View style={styles.changesBadge}>
              <Text style={styles.changesText}>Has Changes</Text>
            </View>
          )}
        </Card>
      </TouchableOpacity>
    );
  };

  if (loading) {
    return (
      <View style={styles.centerContainer}>
        <ActivityIndicator size="large" color={COLORS.primary} />
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.centerContainer}>
        <Text style={styles.errorText}>{error}</Text>
        <TouchableOpacity onPress={loadRevisions} style={styles.retryButton}>
          <Text style={styles.retryButtonText}>Retry</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Revision History</Text>
        {onCompare && selectedVersions.length === 2 && (
          <TouchableOpacity onPress={handleCompare} style={styles.compareButton}>
            <Text style={styles.compareButtonText}>Compare</Text>
          </TouchableOpacity>
        )}
      </View>

      {revisions.length === 0 ? (
        <View style={styles.emptyContainer}>
          <Text style={styles.emptyText}>No revisions available</Text>
        </View>
      ) : (
        <FlatList
          data={revisions}
          renderItem={renderRevision}
          keyExtractor={(item) => item.id}
          contentContainerStyle={styles.listContainer}
        />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
    paddingHorizontal: 16,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    color: COLORS.text.primary,
  },
  listContainer: {
    padding: 16,
    gap: 12,
  },
  revisionCard: {
    marginBottom: 8,
  },
  selectedCard: {
    borderWidth: 2,
    borderColor: COLORS.primary,
  },
  revisionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  versionText: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.text.primary,
  },
  dateText: {
    fontSize: 12,
    color: COLORS.text.secondary,
  },
  revisionDetails: {
    gap: 4,
  },
  detailText: {
    fontSize: 14,
    color: COLORS.text.secondary,
  },
  changesBadge: {
    marginTop: 8,
    alignSelf: 'flex-start',
    backgroundColor: COLORS.primary,
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 4,
  },
  changesText: {
    fontSize: 12,
    color: '#fff',
    fontWeight: '500',
  },
  errorText: {
    fontSize: 14,
    color: COLORS.error,
    textAlign: 'center',
    marginBottom: 16,
  },
  retryButton: {
    backgroundColor: COLORS.primary,
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 8,
  },
  retryButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  compareButton: {
    backgroundColor: COLORS.primary,
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 8,
  },
  compareButtonText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 40,
  },
  emptyText: {
    fontSize: 16,
    color: COLORS.text.secondary,
    textAlign: 'center',
  },
});
