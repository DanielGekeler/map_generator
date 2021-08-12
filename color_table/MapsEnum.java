package com.example.maps;

import net.fabricmc.api.ModInitializer;
import net.minecraft.block.Block;
import net.minecraft.block.BlockState;
import net.minecraft.util.Pair;
import net.minecraft.util.registry.Registry;

import java.io.FileWriter;
import java.io.IOException;
import java.util.*;

public class MapsEnum implements ModInitializer {
    @Override
    public void onInitialize() {
        List<Pair<String, Integer>> l = new ArrayList<>();
        Registry.BLOCK.forEach((Block block) -> {
            HashMap<BlockState, Integer> colors = new HashMap<>();
            block.getStateManager().getStates().forEach((BlockState state) -> {
                int color = state.getMapColor(null, null).id;
                colors.put(state, color);
            });
            Set<Integer> pcolors = new HashSet<>(colors.values());
            if (pcolors.size() == 1) {
                l.add(new Pair<String, Integer>(Registry.BLOCK.getId(block).toString() + "[*]", colors.get(block.getDefaultState())));
            } else {
                colors.forEach((BlockState s, Integer c) -> {
                    l.add(new Pair<String, Integer>(s.toString(), c));
                });
            }
        });
        l.sort(Comparator.comparing(Pair::getRight));
        try {
            FileWriter w = new FileWriter("data.csv");
            l.forEach((Pair<String, Integer> p) -> {
                try {
                    w.write(p.getLeft()+";"+p.getRight()+"\n");
                } catch (IOException e) {
                    e.printStackTrace();
                }
            });
            w.close();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
