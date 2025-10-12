/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */

package com.freshease.backend.modules.permissions;

import java.util.UUID;

import org.springframework.stereotype.Service;

import com.freshease.backend.modules.rolePermissions.RolePermissionRepository;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;

/**
 *
 * @author jemiezler
 */
@Service
@RequiredArgsConstructor
public class PermissionService {
    private final PermissionRepository permissionRepository;
    private final RolePermissionRepository rolePermissionRepository;

    @Transactional
    public PermissionEntity createPermission(PermissionEntity permission) {
        return permissionRepository.save(permission);
    }

    @Transactional()
    public boolean userHasPermission(UUID userId, String permissionName) {
        var perm = permissionRepository.findByName(permissionName).orElse(null);
        if (perm == null)
            return false;
        return rolePermissionRepository.userHasPermission(userId, perm.getId());
    }
}
