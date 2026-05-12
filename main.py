import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
from imblearn.over_sampling import RandomOverSampler
import os
from sklearn.preprocessing import LabelEncoder, normalize
from sklearn.decomposition import PCA
from sklearn.model_selection import train_test_split
from sklearn.svm import SVC
from sklearn.metrics import accuracy_score, confusion_matrix

model = SVC()

def preprocessing(f):
    if not os.path.exists('Images'):
        os.makedirs('Images')

    df = pd.read_csv(f)
    le = LabelEncoder()
    ros = RandomOverSampler(random_state = 0)

    df['type'] = le.fit_transform(df['type'])
    df = df.fillna(df.mean(numeric_only=True))

    y = df['type']
    X = df.drop(columns=['samples', 'type'])
    X_resampled,y_resampled = ros.fit_resample(X,y)
    X_train, X_test, y_train, y_test = train_test_split(X_resampled, y_resampled, test_size=0.2, random_state=42)

    pca = PCA(n_components=0.95)
    X_train_pca = pca.fit_transform(X_train)
    X_test_pca = pca.transform(X_test)

    X_train_pca = normalize(X_train_pca)
    X_test_pca = normalize(X_test_pca)

    return X_train_pca, X_test_pca, y_train, y_test, pca, le

def visualization(X_train_pca, y_train, pca, le, y_test, y_pred):
    plt.figure(figsize=(10, 7))
    scatter = plt.scatter(X_train_pca[:, 0], X_train_pca[:, 1], c=y_train, cmap='viridis', edgecolors='k')
    plt.title('PCA Cluster Analysis')
    plt.xlabel(f'PC1 ({pca.explained_variance_ratio_[0]:.2%} Variance)')
    plt.ylabel(f'PC2 ({pca.explained_variance_ratio_[1]:.2%} Variance)')
    handles, labels = scatter.legend_elements()
    plt.legend(handles, le.classes_, title="Leukemia Types")
    plt.savefig('Images/PCA_Clusters.png')
    plt.show()
    plt.close()

    cm = confusion_matrix(y_test, y_pred)
    plt.figure(figsize=(8, 6))
    sns.heatmap(cm, annot=True, fmt='d', cmap='Blues', xticklabels=le.classes_, yticklabels=le.classes_)
    plt.title('Confusion Matrix')
    plt.ylabel('Actual')
    plt.xlabel('Predicted')
    plt.tight_layout()
    plt.savefig('Images/Confusion_Matrix.png')
    plt.show()
    plt.close()

    plt.figure(figsize=(10, 5))
    plt.bar(range(1, len(pca.explained_variance_ratio_) + 1), pca.explained_variance_ratio_)
    plt.step(range(1, len(pca.explained_variance_ratio_) + 1), np.cumsum(pca.explained_variance_ratio_), where='mid')
    plt.title('Scree Plot')
    plt.xlabel('Principal Components')
    plt.ylabel('Variance')
    plt.savefig('Images/Scree_Plot.png')
    plt.show()
    plt.close()

def train_testing_model(X_train, X_test, y_train, y_test, pca, le):
    model.fit(X_train, y_train)
    y_pred = model.predict(X_test)
    print(f"accuracy: {accuracy_score(y_test, y_pred)}")
    visualization(X_train, y_train, pca, le, y_test, y_pred)

if __name__ == '__main__':
    f = "Leukemia.csv"
    X_train, X_test, y_train, y_test, pca, le = preprocessing(f)
    train_testing_model(X_train, X_test, y_train, y_test, pca, le)
