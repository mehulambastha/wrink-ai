from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class SearchRequest(_message.Message):
    __slots__ = ("topic", "keywords")
    TOPIC_FIELD_NUMBER: _ClassVar[int]
    KEYWORDS_FIELD_NUMBER: _ClassVar[int]
    topic: str
    keywords: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, topic: _Optional[str] = ..., keywords: _Optional[_Iterable[str]] = ...) -> None: ...

class SearchResponse(_message.Message):
    __slots__ = ("search_results_text",)
    SEARCH_RESULTS_TEXT_FIELD_NUMBER: _ClassVar[int]
    search_results_text: str
    def __init__(self, search_results_text: _Optional[str] = ...) -> None: ...
