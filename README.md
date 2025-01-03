# ConfigWeb

ConfigWeb은 실제 인프라의 설정과 추상화된 어플리케이션 패키지의 설정을 매핑하기 위한 도구입니다.

## Summary

ConfigWeb이 해결하고자 하는 바는 간단합니다.  
만약 여러분들이 Redis Cluster를 활용하는 여러가지 패키지를 역할과 책임을 기준으로 분리했다고 가정합시다.  
그러면 여러분들의 어플리케이션의 프로젝트에는 서로 다른 이름의 여러가지 Redis Cluster 클라이언트 역할을 하는 패키지(혹은 클래스)가 존재할 것입니다.  
이 패키지들은 각자의 설정을 가지고 있을 것이고, 이 설정은 실제 연결되는 Redis Cluster의 설정과 일치해야합니다.

이때 각 패키지의 수만큼 Redis Cluster가 존재하면, 각 설정에 각 Redis Cluster의 Addresses와 Password를 입력하면 됩니다.  
하지만 패키지와 실제 구성된 Redis Cluster가 M:N 관계를 가진다면, 이러한 설정은 매우 비효율적입니다.  
그리고 추후 어떤 패키지가 어떤 Redis Cluster 인프라와 연결되어 있는지 파악하기 어려우며, 이를 변경하기도 어렵습니다.  
그렇기에 Redis Cluster의 각 인프라에 대한 설정과 패키지가 어떤 인프라에 연결되는 지 정의하면, 해당 설정에 맞춰서 환경 설정을 매핑해주는 도구가 필요합니다.

이러한 문제를 해결하기 위해 ConfigWeb은 실제 인프라의 설정과 추상화된 어플리케이션 패키지의 설정을 매핑하기 위해 만들어졌습니다.