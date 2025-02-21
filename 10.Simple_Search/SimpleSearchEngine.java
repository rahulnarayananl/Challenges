import java.util.*;

public class SimpleSearchEngine {
    private final Map<String, Set<Integer>> invertedIndex = new HashMap<>();
    private final Map<Integer, String> documents = new HashMap<>();
    private int docIdCounter = 0;

    // Tokenization & Indexing
    public void indexDocument(String content) {
        int docId = docIdCounter++;
        documents.put(docId, content);
        String[] words = content.toLowerCase().split("\\W+"); // Simple tokenizer

        for (String word : words) {
            invertedIndex.computeIfAbsent(word, k -> new HashSet<>()).add(docId);
        }
    }

    // Querying (Simple Boolean Search)
    public List<String> search(String query) {
        String[] words = query.toLowerCase().split("\\W+");
        Set<Integer> resultDocs = new HashSet<>();

        for (String word : words) {
            if (invertedIndex.containsKey(word)) {
                if (resultDocs.isEmpty()) {
                    resultDocs.addAll(invertedIndex.get(word));
                } else {
                    resultDocs.retainAll(invertedIndex.get(word)); // AND operation
                }
            }
        }

        List<String> results = new ArrayList<>();
        for (int docId : resultDocs) {
            results.add(documents.get(docId));
        }
        return results;
    }

    public static void main(String[] args) {
        SimpleSearchEngine searchEngine = new SimpleSearchEngine();
        searchEngine.indexDocument("Lucene is a great search engine");
        searchEngine.indexDocument("Search engines like Solr and Elasticsearch use Lucene");
        searchEngine.indexDocument("Lucene is fast and powerful");

        System.out.println("Search results for 'Lucene search': " + searchEngine.search("Lucene search"));
    }
}