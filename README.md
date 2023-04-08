# Smart Contract Management System

### 하이퍼레저 패브릭에서 스마트 컨트랙트에 대한 재사용성, 확장성을 지원하기 위한 스마트 컨트랙트 관리 시스템입니다.

## 주요 기능
- 스마트 컨트랙트 업로드
- 스마트 컨트랙트 대시보드 - 스마트 컨트랙트 검색 / 비교
- 스마트 컨트랙트 상세 정보 - 설치 / 다운로드 / 트랜잭션


## 구현
|||
|---|---|
|![img7](https://user-images.githubusercontent.com/78259314/230725409-607a57a0-d802-4328-b78e-b2194b9fd61d.png){: width="50" height="50"}|![img6](https://user-images.githubusercontent.com/78259314/230725407-d1db0fb6-fc71-4119-8175-f9b651ae3cd4.png)
|![img8](https://user-images.githubusercontent.com/78259314/230725428-af70880a-5dd2-4c75-99c8-4763ac4e7515.png)||![img9](https://user-images.githubusercontent.com/78259314/230725426-532dad08-5f41-495e-8f3a-3f40a294102d.png)
|![img5](https://user-images.githubusercontent.com/78259314/230725432-1d3bbc23-a9df-4648-bb04-f93578ab3014.png)|

## 평가
> 본 프로젝트에서는 하이퍼레저 패브릭 네트워크와 연결하기 위해 Fabric Gateway SDK를 활용하였으며, 명령어 기반 실행과 SDK 기반 실행 성능을 평가하였다.

테스트 방법
- SDK
  - 플랫폼과 연결하는 시간은 포함되지 않음.
  - REST API 요청 ->응답 사이의 시간을 측정하였음.
  - Jmeter를 활용
- CLI
  - 시간 측정은 서버에서 CLI를 통해 쉘스크립트에 작성된 명령어 set을 실행하는 것으로 소요 시간 측정
  - 반복문-쉘스크립트 -> 동작-쉘스크립트 형태로 반복 수행하였음.
  - 데이터는 파이프라인을 통해 기록하고 이를 쉘스크립트를 통해 min,avg,max 로 종합, 정리함.

특이사항
- CLI 첫 값의 latency가 매우 큼
  - 아마 connection 문제일 듯, Gateway는 connection pool이 존재함.
  - 99% line의 값과 maximum값의 차이가 매우 큼, 이는 한 값이 매우 튐을 추측할 수 있음.
  - => 결과적으로 첫 번째 튀는 값을 제외한 나머지들을 통해서 min,avg,max 값을 평가함.

### Charts
| | |
|---|---|
|![img1](https://user-images.githubusercontent.com/78259314/230723374-26c2b3e4-9c85-409f-94bc-78ec8fea9010.png)|![img2](https://user-images.githubusercontent.com/78259314/230723436-cb8fa374-dc61-417e-9c9c-4d26c184e6b9.png)|
|<p align="center">전체 비교</p>|<p align="center">최대</p>|
|![img3](https://user-images.githubusercontent.com/78259314/230723533-4070e3ba-3ed0-4768-8938-afb6b3928e4c.png)|![img4](https://user-images.githubusercontent.com/78259314/230723537-37b80b56-503f-483a-82cb-57853cca28da.png)|
|<p align="center">평균</p>|<p align="center">최소</p>|










