/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Interface.java to edit this template
 */

package com.freshease.backend.modules.users;

import java.util.UUID;

import org.springframework.stereotype.Service;

import com.freshease.backend.modules.roles.RoleRepository;
import com.freshease.backend.modules.userRoles.UserRoleEntity;
import com.freshease.backend.modules.userRoles.UserRoleRepository;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;

/**
 *
 * @author jemiezler
 */
@Service
@RequiredArgsConstructor
public class UserService {
    private final UserRepository userRepository;
    private final RoleRepository roleRepository;
    private final UserRoleRepository userRoleRepository;

    @Transactional
    public UserEntity createUser(UserEntity user) {
        return userRepository.save(user);
    }

    @Transactional
    public void assignRoleToUser(UUID userId, UUID roleId) {
        if (userRoleRepository.existsByUser_IdAndRole_Id(userId, roleId))
            return;
        var user = userRepository.findById(userId).orElseThrow();
        var role = roleRepository.findById(roleId).orElseThrow();
        userRoleRepository.save(UserRoleEntity.builder().user(user).role(role).build());
    }

}
