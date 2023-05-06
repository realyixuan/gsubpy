class AlphabeticalIterator:
    def __init__(self, collection):
        self._collection = collection
        self._position = 0

    def __next__(self):
        if not self._position < len(self._collection):
            raise StopIteration
            
        value = self._collection[self._position]
        self._position += 1

        return value


class WordsCollection:

    _collection = []

    def __iter__(self):
        return AlphabeticalIterator(self._collection)

    def add_item(self, item):
        self._collection.append(item)


class MyStr(str):
    def isspace(self):
        for c in self:
            if c != ' ':
                return False
        return True


collection = WordsCollection()
collection.add_item(MyStr("First"))
collection.add_item(MyStr(" "))
collection.add_item(MyStr("Second"))
collection.add_item(MyStr("Third"))

print("Straight traversal:")
for item in collection:
    if item.isspace() is False:
        print(item)

