package com.freshease.backend.modules.userRoles;

import java.util.UUID;

import org.springframework.stereotype.Service;

import com.freshease.backend.modules.roles.RoleRepository;
import com.freshease.backend.modules.users.UserRepository;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class UserRoleService {
    private final UserRoleRepository userRolesRepository;
    private final UserRepository usersRepository;
    private final RoleRepository rolesRepository;

    @Transactional
    public void assignRoleToUser(UUID userId, UUID roleId) {
        if (userRolesRepository.existsByUser_IdAndRole_Id(userId, roleId))
            return;
        var user = usersRepository.findById(userId).orElseThrow();
        var role = rolesRepository.findById(roleId).orElseThrow();
        userRolesRepository.save(UserRoleEntity.builder().user(user).role(role).build());
    }
}
